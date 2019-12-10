const $ = require('jquery');
const fs = require('fs');
const readline = require('readline');
// This will wait for the astilectron namespace to be ready
document.addEventListener('astilectron-ready', function() {
    // This will send a message to GO
    // astilectron.sendMessage("hello", function(message) {
    //     console.log("received " + message)
    // });

    // This will listen to messages sent by GO
    astilectron.onMessage(function(message) {
        // Process message
        if (message === "hello") {
            return "received"
        }
    });
});

$("#addProxy").click(
    function () {
        astilectron.sendMessage("addProxy://name-socks5://127.0.0.1:1080\naddProxy://name-http://127.0.0.1:8080",function (message) {
            console.log("received: "+message);
            alert(message)
        });
    }
);

$("#deleteProxy").click(
    function () {
        astilectron.sendMessage("deleteProxy://name",function (message) {
          console.log("received: "+message);
          alert(message)
        })
    }
);

function View(view) {
    if (view === "proxy") {
        proxyTableInit()
    }else if (view === "rule"){
        ruleTableInit()
    }else if(view === "setting"){
        settingInit()
    }

    $("#content").children().css("display","none");
    $("#"+view).css("display","block");
    $("#collapsibleNavbar").collapse('hide');
}

function ruleTableInit(){
    const rl = readline.createInterface({
        input: fs.createReadStream('./rule/rule.config')
    });
    $("#ruleTable").empty();
    rl.on('line', (line) => {
        const arr = line.split(' ');
        $("#ruleTable").append(`
                        <tr>
                            <td>${arr[0]}</td>
                            <td>${arr[1]}</td>
                            <td><button type="button" class="btn btn-light" onclick="deleteRule('${arr[0]}')">DELETE</button></td>
                        </tr>`)
        // console.log('访问时间：%s %s', arr[0], arr[1]);
    });


    fs.readFile('./config/config.json',function (err,data) {
        if (err) {
            console.log(err)
        }
        let proxy = JSON.parse(data);
        let nodes = proxy["nodes"];
        let addRuleProxy = $("#addRuleProxy");
        addRuleProxy.empty();
        for (const key in nodes) {
            if (nodes.hasOwnProperty(key)) {
                addRuleProxy.append(`<option>${key}</option>`)
            }
        }
    })
}


function settingInit() {
    let moreSetting = $("#moreSetting");
    let bypassSetting = $("#bypassSetting");
    let directSetting = $("#directSetting");
    let proxySetting = $("#proxyOnlySetting");

    fs.readFile('./config/config.json',function (err,data) {
        if(err){
            console.log(err)
        }
        let proxy = JSON.parse(data);
        let setting = proxy["setting"];
        let nodes = proxy["nodes"];
        $("#dnsSetting").val(setting["dns"]);
        $("#socks5Setting").val(setting["socks_5"]);
        $("#httpSetting").val(setting["http"]);
        bypassSetting.prop("checked",setting["bypass"]);
        directSetting.prop("checked",setting["direct"]);
        if (!bypassSetting.prop("checked")){
            moreSetting.css("display","block")
        }else{
            moreSetting.css("display","none")
        }
        if (directSetting.prop("checked")){
            proxySetting.attr("disabled","disabled")
        }else{
            proxySetting.removeAttr("disabled")
        }
        proxySetting.empty();
        for (const key in nodes){
            if (nodes.hasOwnProperty(key)){
                proxySetting.append(`<option>${key}</option>`)
            }
        }
        if (setting["proxy"] in nodes){
            proxySetting.val(setting["proxy"]);
        }
    });
    bypassSetting.change(
        function () {
            if ($(this).prop("checked")){
                moreSetting.css("display","none")
            }else {
                moreSetting.css("display","block")
            }
        }
    );
    directSetting.change(
        function () {
            if ($(this).prop("checked")){
                proxySetting.attr("disabled","disabled")
            }else {
                proxySetting.removeAttr("disabled")
            }
        }
    )
}

function getAlert(str) {
    return `
<div class="alert alert-success alert-dismissible fade show">
  <button type="button" class="close" data-dismiss="alert">&times;</button>
  <strong>${str}</strong>
</div>`
}
const deleteProxy = id => {
    $("#proxyWar").html(getAlert("删除"+id+"成功!"));
    proxyTableInit();
};

const deleteRule = id => {
    $("#ruleWar").html(getAlert("删除"+id+"成功!"));
    ruleTableInit();
};

function proxyTableInit(){
    $("#proxyTable").empty();
    fs.readFile('./config/config.json',function (err,data) {
        if(err){
            console.log(err)
        }
        let proxy = JSON.parse(data);
        let nodes = proxy["nodes"];
        for (const key in nodes){
            if (nodes.hasOwnProperty(key)){
                let node = nodes[key];
                $("#proxyTable").append(`
                        <tr>
                            <td>${key}</td>
                            <td>${node["Scheme"]}</td>
                            <td>${node["Host"]}</td>
                            <td><button type="button" class="btn btn-light" onclick="deleteProxy('${key}')">DELETE</button></td>
                        </tr>`)
            }
        }
    })
}

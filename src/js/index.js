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
    $("#navi").collapse('hide');
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
                            <td><button type="button" class="btn bth-default">DELETE</button></td>
                        </tr>`)
        // console.log('访问时间：%s %s', arr[0], arr[1]);
    });
}


function settingInit() {
    let moreSetting = $("#moreSetting");
    let bypassSetting = $("#bypassSetting");
    let directSetting = $("#directSetting");
    let proxySetting = $("#proxyOnlySetting");
    if (bypassSetting.prop("checked")){
        moreSetting.css("display","none")
    }
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
                            <td>${node["Scheme"]}</td>
                            <td>${node["Host"]}</td>
                        </tr>`)
            }
        }
    })
}

$(document).ready(
    function () {
    }
);

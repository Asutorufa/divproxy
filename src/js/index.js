const $ = require('jquery');
const fs = require('fs');
const readline = require('readline');
// This will wait for the astilectron namespace to be ready
document.addEventListener('astilectron-ready', function() {
    // This will send a message to GO
    // astilectron.sendMessage("hello", function(message) {
    //     console.log("received " + message)
    // });

});



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
    let proxySetting = $("#proxyOnlySetting");
    let proxyModeSetting = $("#proxyModeSetting");
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
        if (setting["bypass"]){
            proxyModeSetting.val("BYPASS")
        }else{
            if (setting["direct"]){
                proxyModeSetting.val("DIRECT")
            }else{
                proxyModeSetting.val("PROXY")
            }
        }
        if (proxyModeSetting.val() !== "PROXY"){
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
    proxyModeSetting.change(
        function () {
            if ($(this).val() !== "PROXY"){
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

$("document").ready(
    function () {
        $('[data-toggle="popover"]').popover()
    }
);
// const path = require("path");
// console.log(path.resolve("./src/connectBanckend.js"));
// const {startProxy} = require(path.resolve("./src/js/connectBanckend.js"));

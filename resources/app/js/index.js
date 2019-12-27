const $ = require('jquery');
const fs = require('fs');
const readline = require('readline');


function getAlert(str) {
    return `
<div class="alert alert-success alert-dismissible fade show">
  <button type="button" class="close" data-dismiss="alert">&times;</button>
  <strong>${str}</strong>
</div>`
}

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


$("document").ready(
    function () {
        $('[data-toggle="popover"]').popover();
    }
);

$("#refreshSettingButton").click(
    function () {
        settingInit();
        $("#settingWar").html(getAlert("refresh setting!"));
    }
);

function editProxyModalShow(name,scheme,host) {
    $("#editProxyName").val(name);
    $("#editProxyScheme").val(scheme.toUpperCase());
    $("#editProxyHost").val(host);
    $("#myModalEditProxy").modal('show')
}

/*
            Init Start
 */
function ruleTableInit(){
    const rl = readline.createInterface({
        input: fs.createReadStream('./resources/app/config/rule.config')
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


    fs.readFile('./resources/app/config/config.json',function (err,data) {
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
        addRuleProxy.append(`<option>block</option>`)
    })
}


function settingInit() {
    let proxySetting = $("#proxyOnlySetting");
    let proxyModeSetting = $("#proxyModeSetting");
    fs.readFile('./resources/app/config/config.json',function (err,data) {
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

function proxyTableInit(){
    $("#proxyTable").empty();
    fs.readFile('./resources/app/config/config.json',function (err,data) {
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
                            <td>${node["Scheme"].toUpperCase()}</td>
                            <td>${node["Host"]}</td>
                            <td>
                                <button type="button" class="btn btn-light" onclick="deleteProxy('${key}')">DELETE</button>
                                <button type="button" class="btn btn-light" onclick="editProxyModalShow('${key}','${node["Scheme"].toUpperCase()}','${node["Host"]}')">EDIT</button>
                            </td>
                        </tr>`)
            }
        }
    })
}
/*
            Init end
*/
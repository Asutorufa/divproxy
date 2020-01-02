const moment = require("moment");
function startProxy() {
    astilectron.sendMessage("startProxy://", function (message) {
        console.log("received: " + message);
        $("#homepageWar").html(getAlert(message))
    })
}

function stopProxy() {
    astilectron.sendMessage("stopProxy://", function (message) {
        console.log("received: " + message);
        $("#homepageWar").html(getAlert(message))
    })
}

function restartProxy() {
    astilectron.sendMessage("restartProxy://", function (message) {
        console.log("received: " + message);
        $("#homepageWar").html(getAlert(message))
    })
}

document.addEventListener('astilectron-ready', function() {
    // This will listen to messages sent by GO
    astilectron.onMessage(function(message) {
        // Process message
        // console.log(message);
        $("#logCode").prepend("<span style='color: #adb7c3'>"+moment().format('MMMM Do, hh:mm:ss a')+"</span>    "+message+"\n")
    });

    $("#startProxy").click(
        function(){
            startProxy()
        }
    );
    $("#restartProxy").click(
        function(){
            restartProxy()
        }
    );
    $("#stopProxy").click(
        function(){
            stopProxy()
        }
    );


    $("#editProxyButton").click(
        function () {
            let editProxyHost = $("#editProxyHost");
            let name = $("#editProxyName").val();
            let scheme = $("#editProxyScheme").val();
            let host = editProxyHost.val();
            if (!/^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}$/.test(host)){
                editProxyHost.popover('hide').popover('show');
                return
            }
            astilectron.sendMessage("addProxy://"+name+"-"+scheme+"://"+host,function (message) {
                $("#ruleWar").html(getAlert(message));
                restartProxy();
                proxyTableInit();
                $("#myModalEditProxy").modal('hide')
            });
        });




    $("#addProxyButton").click(
        function () {
            let addProxyName = $("#addProxyName");
            let addProxyHost = $("#addProxyHost");
            let name = addProxyName.val();
            let scheme = $("#addProxyScheme").val().toLowerCase();
            let Host = addProxyHost.val();
            if (name === "") {
                addProxyName.popover('hide').popover('show');
                return
            }
            if (!/^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}$/.test(Host)){
                addProxyHost.popover('hide').popover('show');
                return
            }
            console.log(name+scheme+Host);
            $('#myModalProxy').modal('hide');
            astilectron.sendMessage("addProxy://"+name+"-"+scheme+"://"+Host,function (message) {
                // console.log("received: "+message);
                $("#proxyWar").html(getAlert(message));
                proxyTableInit();
            });
        }
    );

    $("#addRuleButton").click(
        function(){
            let addRuleRule = $("#addRuleRule");
            if (addRuleRule.val()===""){
                addRuleRule.popover('hide').popover('show');
                return
            }
            let rule = addRuleRule.val();
            let proxy = $("#addRuleProxy").val();
            if (!/(^(?=^.{3,255}$)[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$|^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\/\d{1,2}$)/.test(rule)){
                addRuleRule.popover('hide').popover('show');
                return;
            }
            astilectron.sendMessage("addRule://"+rule+"-"+proxy, function (message) {
                console.log("received: " + message);
                $("#ruleWar").html(getAlert(message));
                ruleTableInit();
                $("#myModalRule").modal('hide');
                restartProxy()
            });
        });


    $("#applySettingButton").click(
        function () {
            console.log("apply setting!");
            let DNS = $("#dnsSetting").val();
            let socks5 = $("#socks5Setting").val();
            let http = $("#httpSetting").val();
            let proxyMode = $("#proxyModeSetting").val();
            let onlyProxy = $("#proxyOnlySetting").val();
            // console.log(DNS,socks5,http,proxyMode,onlyProxy);
            astilectron.sendMessage("applySetting://DNS-"+DNS+":socks5-"+socks5+":http-"+http+":proxyMode-"+proxyMode+":onlyProxy-"+onlyProxy,function (message) {
                console.log("received: "+message);
                $("#settingWar").html(getAlert(message));
                proxyTableInit();
                restartProxy()
            });
        }
    );


    $("#refreshSettingButton").click(
        function () {
            settingInit();
            $("#settingWar").html(getAlert("refresh setting!"));
        }
    );

});



const deleteProxy = id => {
    astilectron.sendMessage("deleteProxy://"+id,function (message) {
        console.log("received: "+message);
        $("#proxyWar").html(getAlert(message));
        restartProxy();
        proxyTableInit();
    });
};

const deleteRule = id => {
    astilectron.sendMessage("deleteRule://"+id,function (message) {
        $("#ruleWar").html(getAlert(message));
        restartProxy();
        ruleTableInit();
    });
};

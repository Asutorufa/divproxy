document.addEventListener('astilectron-ready', function() {
    // This will listen to messages sent by GO
    astilectron.onMessage(function(message) {
        // Process message
        // console.log(message);
        $("#logCode").prepend(new Date()+"    "+message+"\n")
    });

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
                $("#ruleWar").html(getAlert(message))
                ruleTableInit();
                $("#myModalRule").modal('hide')
            });
        });

});

const deleteProxy = id => {
    astilectron.sendMessage("deleteProxy://"+id,function (message) {
        console.log("received: "+message);
        $("#proxyWar").html(getAlert(message));
        proxyTableInit();
    });
};


const deleteRule = id => {
    astilectron.sendMessage("deleteRule://"+id,function (message) {
        $("#ruleWar").html(getAlert(message));
        ruleTableInit();
    });
};
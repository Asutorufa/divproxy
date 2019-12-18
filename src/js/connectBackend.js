document.addEventListener('astilectron-ready', function() {
    // This will listen to messages sent by GO
    astilectron.onMessage(function(message) {
        // Process message
        // console.log(message);
        if (message.startsWith("Mode: ")){
            const date = new Date();
            $("#logCode").prepend(date+"    "+message+"\n")
        }
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
                // alert("name is empty!");
                $("#addProxyName").popover('hide').popover('show');
                return
            }
            if (!/^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}$/.test(Host)){
                // alert("Host format error!");
                $("#addProxyHost").popover('hide').popover('show');
                return
            }else {
                $("#addProxyHost").popover('dispose')
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

});

const deleteProxy = id => {
    astilectron.sendMessage("deleteProxy://"+id,function (message) {
        console.log("received: "+message);
        $("#proxyWar").html(getAlert(message));
        proxyTableInit();
    });
};
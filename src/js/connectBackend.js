
document.addEventListener('astilectron-ready', function() {

    // This will listen to messages sent by GO
    astilectron.onMessage(function(message) {
        // Process message
        // console.log(message);
        if (message.startsWith("Mode: ")){
            $("#logCode").prepend(message+"\n")
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
});

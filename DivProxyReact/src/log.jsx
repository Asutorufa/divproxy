import React from 'react';


function Log() {
    return(
        <div id={'log'} style={{display:'none',position:'absolute',height:'85%',width:'90%',overflow:'auto',background: '#f3f6f9',color: '#73808f',padding: '10px',fontSize: '16px'}}>
            <pre id="logCode">

            </pre>
        </div>
    )
}

export default Log;
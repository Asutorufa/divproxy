import React from 'react';

const buttons = [
    <button type="button" className="btn btn-secondary" id="startProxy">start</button>,
    <button type="button" className="btn btn-secondary" id="restartProxy">restart</button>,
    <button type="button" className="btn btn-secondary" id="stopProxy">stop</button>
];

function Homepage() {
    return(
        <div id={"homepage"}>
            <div id={'homepageWar'}>

            </div>
            {buttons}
        </div>
    )
}

export default Homepage;
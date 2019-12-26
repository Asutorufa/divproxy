import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import Nav from "./nav";
import Homepage from "./homepage";
import Log from "./log";
import Rule from "./rule";
import Setting from "./setting";
import Proxy from "./proxy";
import * as serviceWorker from './serviceWorker';

function Content(){
    return(
        <div id={"content"} style={{padding:'0 70px 0 70px'}}>
            <Homepage/>
            <Log/>
            <Proxy/>
            <Setting/>
            <Rule/>
        </div>
    );
}


// ReactDOM.render(<App />, document.getElementById('root'));
// ReactDOM.render(<Content/>,document.getElementById('root'));
var headElement =[
    <script>require('popper.js')</script>,
    <script>require('bootstrap');</script>,
    <script>window.$ = window.jQuery = require('jquery');</script>,
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossOrigin="anonymous"/>];
ReactDOM.render(headElement,document.getElementsByTagName('head')[0]);
ReactDOM.render([<Nav/>,<Content/>, <script src='js/homepage.js'/>, <script src='js/connectBackend.js'/>],document.getElementsByTagName("body")[0]);
// ReactDOM.render([<Homepage/>,<Log/>,<Proxy/>,<Setting/>,<Rule/>],document.getElementById('content'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();

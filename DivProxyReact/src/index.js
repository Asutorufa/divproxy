import React from 'react';
import ReactDOM from 'react-dom';
import Nav from "./nav";
import Homepage from "./homepage";
import Log from "./log";
import Rule from "./rule";
import Setting from "./setting";
import Proxy from "./proxy";
import * as serviceWorker from './serviceWorker';

function Body(){
    return(
        <body>
        <Nav/>
        <div id={"content"} style={{padding: '70px'}}>
            <Homepage/>
            <Log/>
            <Proxy/>
            <Setting/>
            <Rule/>
        </div>
        <script src="../app/js/index.js"/>
        <script src="../app/js/connectBackend.js"/>
        </body>
    );
}

function Heads() {
    return(
        <head>
            <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossOrigin="anonymous"/>
            <style>
                {`
                    /*- scrollbar -*/
                    ::-webkit-scrollbar {
                    width: 5px;
                    height: 5px;
                }
                    ::-webkit-scrollbar-thumb{
                    background-color: #999;
                    -webkit-border-radius: 5px;
                    border-radius: 5px;
                }
                    ::-webkit-scrollbar-thumb:vertical:hover{
                    background-color: #666;
                }
                    ::-webkit-scrollbar-thumb:vertical:active{
                    background-color: #333;
                }
                    ::-webkit-scrollbar-button{
                    display: none;
                }
                    ::-webkit-scrollbar-track{
                    background-color: #f1f1f1;
                }`
                }
            </style>
            <title>DivProxy</title>

            <script>require('popper.js');</script>
            <script>require('bootstrap');</script>
            <script>window.$ = window.jQuery = require('jquery');</script>
        </head>
    )
}


ReactDOM.render([<Heads/>,<Body/>],document.getElementsByTagName('html')[0]);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();

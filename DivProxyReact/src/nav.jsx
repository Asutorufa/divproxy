import React from 'react';
function Nav() {
    return(
        <nav className="navbar navbar-expand-sm bg-light navbar-light fixed-top" role="navigation">
            <a className="navbar-brand" href="javascript: View('homepage')">DivProxy</a>
            <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#collapsibleNavbar">
                <span className="navbar-toggler-icon"/>
            </button>
            <div className="collapse navbar-collapse" id="collapsibleNavbar">
                <ul className="navbar-nav">
                    <li className="nav-item"><a className="nav-link" href="javascript: View('log')">LOG</a></li>
                    <li className="nav-item"><a className="nav-link" href="javascript: View('proxy')">PROXY</a></li>
                    <li className="nav-item"><a className="nav-link" href="javascript: View('rule')">RULE</a></li>
                    <li className="nav-item"><a className="nav-link" href="javascript: View('setting')">SETTING</a></li>
                </ul>
            </div>
        </nav>
    )
}

export default Nav;
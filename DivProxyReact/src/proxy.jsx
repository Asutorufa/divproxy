import React from 'react';

function proxyTable(){
    return (
        <table className="table">
            <tbody>
            <tr>
                <th>NAME</th>
                <th>SCHEME</th>
                <th>HOST</th>
                <th>OPTIONS</th>
            </tr>
            </tbody>
            <tbody id="proxyTable">

            </tbody>
        </table>
            )
}

function proxyModal(){
    return(
        <div className="modal fade" id="myModalProxy">
            <div className="modal-dialog">
                <div className="modal-content">

                    <div className="modal-header">
                        <h4 className="modal-title">ADD PROXY</h4>
                        <button type="button" className="close" data-dismiss="modal">&times;</button>
                    </div>

                    <div className="modal-body">
                        <div className="input-group mb-3">
                            <div className="input-group-prepend">
                                <span className="input-group-text">NAME</span>
                            </div>
                            <label htmlFor="addProxyName"/>
                            <input type="text" className="form-control"
                                   data-toggle="popover" data-trigger="focus"
                                   data-placement="top"
                                   data-animation="true"
                                   data-content="Name Can Not Empty"
                                   id="addProxyName" placeholder="Name.."/>
                        </div>
                        <div className="input-group mb-3">
                            <div className="input-group-prepend">
                                <span className="input-group-text">SCHEME</span>
                            </div>
                            <label htmlFor="addProxyScheme"/>
                            <select id="addProxyScheme" className="custom-select" data-style="btn-primary">
                                <option>SOCKS5</option>
                                <option>HTTP</option>
                            </select>
                        </div>
                        <div className="input-group mb-3">
                            <div className="input-group-prepend">
                                <span className="input-group-text">HOST</span>
                            </div>
                            <label htmlFor="addProxyHost"/>
                            <input type="text" className="form-control"
                                   data-toggle="popover" data-trigger="focus"
                                   data-placement="top"
                                   data-animation="true"
                                   data-html="true"
                                   data-content="Format Like This:<br/><code>127.0.0.1:8080?username=xx&password=xxx</code>"
                                   id="addProxyHost" placeholder="127.0.0.1:1080.."/>
                        </div>
                    </div>

                    <div className="modal-footer">
                        <button type="button" id="addProxyButton" className="btn btn-primary">添加</button>
                        <button type="button" className="btn btn-secondary" data-dismiss="modal">取消</button>
                    </div>

                </div>
            </div>
        </div>
    )
}

function proxyEditModal() {
    return(
        <div className="modal fade" id="myModalEditProxy">
            <div className="modal-dialog">
                <div className="modal-content">

                    <div className="modal-header">
                        <h4 className="modal-title">EDIT PROXY</h4>
                        <button type="button" className="close" data-dismiss="modal">&times;</button>
                    </div>

                    <div className="modal-body">
                        <div className="input-group mb-3">
                            <div className="input-group-prepend">
                                <span className="input-group-text">NAME</span>
                            </div>
                            <label htmlFor="editProxyName"/>
                            <input type="text" className="form-control" disabled id="editProxyName" value="test"/>
                        </div>
                        <div className="input-group mb-3">
                            <div className="input-group-prepend">
                                <span className="input-group-text">SCHEME</span>
                            </div>
                            <label htmlFor="editProxyScheme"/>
                            <select id="editProxyScheme" className="custom-select" data-style="btn-primary">
                                <option>SOCKS5</option>
                                <option>HTTP</option>
                            </select>
                        </div>
                        <div className="input-group mb-3">
                            <div className="input-group-prepend">
                                <span className="input-group-text">HOST</span>
                            </div>
                            <label htmlFor="editProxyHost"/>
                            <input type="text" className="form-control"
                                   data-toggle="popover" data-trigger="focus"
                                   data-placement="top"
                                   data-animation="true"
                                   data-html="true"
                                   data-content="Format Like This:<br/><code>127.0.0.1:8080?username=xx&password=xxx</code>"
                                   id="editProxyHost" placeholder="127.0.0.1:1080.."/>
                        </div>
                    </div>

                    <div className="modal-footer">
                        <button type="button" id="editProxyButton" className="btn btn-primary">确认修改</button>
                        <button type="button" className="btn btn-secondary" data-dismiss="modal">取消</button>
                    </div>

                </div>
            </div>
        </div>
    )
}

function Proxy() {
    return(
        <div id={'proxy'} style={{display:'none'}}>
            <div id={'proxyWar'}>

            </div>
            {proxyTable()}
            <button className="btn btn-primary" data-toggle="modal" data-target="#myModalProxy">addProxy</button>
            {proxyModal()}
            {proxyEditModal()}
        </div>
    )
}
export default Proxy;

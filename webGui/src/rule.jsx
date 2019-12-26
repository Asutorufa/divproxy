import React from 'react';

function ruleModal() {
    return(
        <div className="modal fade" id="myModalRule">
            <div className="modal-dialog">
                <div className="modal-content">

                    <div className="modal-header">
                        <h4 className="modal-title">ADD RULE</h4>
                        <button type="button" className="close" data-dismiss="modal">&times;</button>
                    </div>

                    <div className="modal-body">
                        <div className="input-group mb-3">
                            <div className="input-group-prepend">
                                <span className="input-group-text">CIDR/DOMAIN</span>
                            </div>
                            <label htmlFor="addRuleRule"/>
                            <input type="text" className="form-control" id="addRuleRule"
                                   data-toggle="popover" data-trigger="focus"
                                   data-placement="top"
                                   data-animation="true"
                                   data-html="true"
                                   data-content="DOMAIN Format: <code>google.com</code><br/>CIDR Format:<code>10.2.2.1/18</code>"
                                   placeholder="10.2.2.1/18.."/>
                        </div>
                        <div className="input-group mb-3">
                            <div className="input-group-prepend">
                                <span className="input-group-text">PROXY</span>
                            </div>
                            <label htmlFor="addRuleProxy"/>
                            <select id="addRuleProxy" className="custom-select" data-style="btn-primary">

                            </select>
                        </div>
                    </div>

                    <div className="modal-footer">
                        <button type="button" className="btn btn-primary" id="addRuleButton">添加</button>
                        <button type="button" className="btn btn-secondary" data-dismiss="modal">取消</button>
                    </div>

                </div>
            </div>
        </div>
    )
}
function ruleTable() {
    return(
        <table className="table">
            <tbody>
            <tr>
                <th>RULE</th>
                <th>PROXY</th>
                <th>OPTION</th>
            </tr>
            </tbody>
            <tbody id="ruleTable">

            </tbody>
        </table>
    )
}
function Rule() {
    return(
        <div id={'rule'} style={{display:'none'}}>
            {ruleTable()}
            <button className="btn btn-primary" data-toggle="modal" data-target="#myModalRule">ADD RULE</button>
            {ruleModal()}
        </div>
    )
}

export default Rule;
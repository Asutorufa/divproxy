import React from 'react';

function settingInputText(id,name) {
    return(
            <div className="input-group mb-3">
                <div className="input-group-prepend">
                    <span className="input-group-text">{name}</span>
                </div>
                <label htmlFor={id}/>
                <input type="text" className="form-control" id={id}/>
            </div>
    )
}

function settingInputSelect(id,name,option) {
    return(
        <div className="input-group mb-3">
            <div className="input-group-prepend">
                <span className="input-group-text">{name}</span>
            </div>
            <label htmlFor={id}/>
            <select id={id} className="custom-select" data-style="btn-primary">
                {option}
            </select>
        </div>
    )
}
function settingFrom() {
    return(
        // eslint-disable-next-line jsx-a11y/no-redundant-roles
        <form className="form" role="form">
            {settingInputText('dnsSetting','DNS')}
            {settingInputText('socks5Setting','SOCKS5')}
            {settingInputText('httpSetting',"HTTP")}
            {settingInputSelect('proxyModeSetting','PROXY MODE',[<option>BYPASS</option>, <option>DIRECT</option>, <option>PROXY</option>])}
            {settingInputSelect('proxyOnlySetting','PROXY',[])}
        </form>
    )
}

function Setting() {
    return(
        <div id={'setting'} style={{display:'none'}}>
            <div id={'settingWar'}>

            </div>
            {settingFrom()}
            <button type="button" className="btn btn-primary" id="applySettingButton">Apply</button>
            <button type="button" className="btn btn-secondary" id="refreshSettingButton">Refresh</button>
        </div>
    )
}

export default Setting;
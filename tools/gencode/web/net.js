(function() {
var config = {
    url: "http://127.0.0.1:5000",
    version: "1.0",
    channel: "mine",
    heartTime: 10000,
}
var exports = {};
window.net = exports;
exports.config = config;
exports.time = new proto.ServerTime();
exports.login_id = 0;
var genGetUrl= function(url){
    url += "?";
    for(var i = 1; i < arguments.length; i++) {
        var args = arguments[i]
        for(var k in args){
            url += k + "=" + ((args[k] == void 0) ? "" : args[k]) + "&";
        }
    }
    return url.substr(0, url.length - 1);
}
var _EventDispatch = function(){
    this._listeners = [];
    this.on = function(callback){
        this._listeners.push(callback);
    }
    this.off = function(callback){
        for(var j = 0,len = this._listeners.length; j < len; j++){
            if(this._listeners[j] == callback){
                this._listeners.splice(j, 1);
                break;
            }
        }
    }
    this.clear = function() {
        this._listeners = []
    }
}
var _httpGet = function(url, callback){
    var xhr = new XMLHttpRequest();
    var json = {error: 45, errstr: "http net failed"};
    xhr.onreadystatechange = function (){
        var ok = xhr.readyState == 4 && (xhr.status >= 200 && xhr.status < 400);
        if(ok) {
            console.log(config.url + url, xhr.responseText);
            json = JSON.parse(xhr.responseText);
            ok = json.error == 0;
            if(callback != void 0) {
                callback(ok, json, xhr);
            }
        } else if(xhr.readyState == 4) {
            callback(false, json);
        }
    };
    xhr.open("GET", config.url + url, true);
    xhr.send();
}
var C2S = function(name, args) {
    var c2s = new proto.C2S();
    var fname = name.substring(0,1).toUpperCase()+name.substring(1).toLowerCase();
    var pb = new proto[name]();    
    pb.setId(exports.login_id);
    c2s["set" + fname](pb);
    pb.send = function () {
        exports.logic._ws.send(c2s.serializeBinary());
    }
    if(args == void 0) {
        return pb;
    } else {
        for(var k in args){
            pb["set" + k.substring(0,1).toUpperCase()+k.substring(1).toLowerCase()](args[k]);
        }
        pb.send();
    }
}
var _dispatchS2C = function(buf){
    var s2c = proto.S2C.deserializeBinary(buf);
    var key = s2c.getKey();
    var ikey = "";
    if(key == ""){
        for(var k in exports.logic){
            ikey = k;
            k = k.substring(0,1).toUpperCase()+k.substring(1).toLowerCase();
            if(typeof s2c["has"+k] == "function" && s2c["has"+k]()){
                key = k;
                break;
            }
        }
    }
    if(key && key != "") {
        var pbs2c = s2c["get"+key]();
        if(pbs2c.getError() == 0) {
            if(key == "Gamerlogins2c"){
                exports.logic._session = pbs2c.getMain().getSession();
            } else if(key == "Servertimes2c" || s2c.key == "GamerLoginGetDataS2C"){
                exports.time = pbs2c.getTime();
            } else if(key == "GamerNotifyLoginOtherS2C"){
                exports.logic._session = "";
            }
        }
        var dispatch = exports.logic[ikey];
        if(dispatch !== void 0 && dispatch != null){
            for(var j = 0,len = dispatch._listeners.length; j < len; j++){
                dispatch._listeners[j](pbs2c);
            }
        }
    }
    if(s2c.getError() > 0){
        for(var j = 0,len = exports.logic.onError._listeners.length; j < len; j++){
            exports.logic.onError._listeners[j](s2c.getError());
        }
    }
}
exports.auth = {
    config: function(callback){
        _httpGet(genGetUrl("/config"), callback);
    },
    register: function(name, passwd, callback){
        _httpGet(genGetUrl("/register", {channel:"mine", name:name, passwd:passwd}), callback);
    },
    login: function(name, passwd, callback) {
        _httpGet(genGetUrl("/login", {channel:config.channel, name:name, passwd:passwd}), callback);
    },
    newRole: function(session, name, type, server, callback) {
        _httpGet(genGetUrl("/newrole", {session:session, name:name, type:type, server:server}), callback);
    },
    useRole: function(session, id, callback) {
        exports.login_id = id;
        exports.logic._session = "";
        _httpGet(genGetUrl("/userole", {session:session, id:id}), function(ok, json, xhr){
            if(ok){ 
                if(config.url.indexOf("http://") == 0) {
                    exports.logic._addr = config.url.replace("http:", "ws:") + "/proxy";
                } else if(config.url.indexOf("https://") == 0) {
                    exports.logic._addr = config.url.replace("https:", "wss:") + "/proxy";
                }
            }
            callback(ok, json, xhr);
        });
    },
    newAndUseRole: function(session, name, type, server, callback){
        var self = this;
        this.newRole(session, name, type, server, function(ok, json){
            ok ? self.useRole(json.session, json.roles[0].id, callback) : callback(ok, json);
        })
    },
    mustLogin: function(name, passwd, roleName, type, server, callback){
        var self = this;
        self.login(name, passwd, function(ok, json){
            ok ? (json.roles == null ? self.newAndUseRole(json.session, roleName, type, server, callback) : self.useRole(json.session, json.roles[0].id, callback)) :
            json.error == 45 ? callback(ok, json) : self.register(name, passwd, function(ok, json) {
                ok ? self.newAndUseRole(json.session, name, type, server, callback) : callback(ok, json);
            });
        });
    }
}

exports.logic = {
_addr: "",
_ws: null,
_session:"",
_reconnCnt:0,
_heart:null,
connect: function(addr){
    this._reconnCnt++;
    if(addr != void 0) {
        this._addr = addr;
    }
    console.log("connect to addr:", this._addr, " session:", this._session);
    if(this._ws != null){
        this._ws.onerror = function(){}
        this._ws.onclose = function(){}
        this._ws.close();
        this._ws = null;
    }
    var onopen = function(){
        this._reconnCnt = 0;
        clearInterval(this._heart);
        this._heart = setInterval(function(){
            this.serverTimeC2S();
        }.bind(this), config.heartTime);
        this.gamerLoginC2S(this._session);
        console.log("connect to addr:", this._addr, " ok session:", this._session);
        for(var j = 0,len = this.onConnect._listeners.length; j < len; j++){
            this.onConnect._listeners[j]();
        }
    }.bind(this);
    var onmessage = function(e){
        if(e.data instanceof ArrayBuffer){
            _dispatchS2C(new Uint8Array(e.data, e.data.byteOffset, e.data.byteLength));
        } else if (e.data instanceof Blob) {
            var reader = new FileReader();
            reader.readAsArrayBuffer(e.data);
            reader.onload = function(evt){
                if(evt.target.readyState == FileReader.DONE){
                    _dispatchS2C(new Uint8Array(evt.target.result));
                }
            }
        }
    };
    var onclose = function(e){
        console.log("onclose", e);
        clearInterval(this._heart);
        console.log("reconnect to addr:", this._addr, " session:", this._session);
        if(this._addr == "" || this._session == ""){
            clearInterval(this._heart);
            for(var j = 0,len = this.onClose._listeners.length; j < len; j++){
                this.onClose._listeners[j]();
            }
        } else {
            this.connect();
            this._heart = setInterval(function(){
                this.connect();
            }.bind(this), 1000);
        }
    }.bind(this);
    var onerror = function(e){
        console.log("onerror", e);
        onclose();
    }.bind(this);

    this._ws = new WebSocket(this._addr);
    this._ws.binaryType="arraybuffer";
    this._ws.onopen = onopen;
    this._ws.onmessage = onmessage;
    this._ws.onerror = onerror;
    this._ws.onclose = onclose;
},
onError: new _EventDispatch(),
onConnect:  new _EventDispatch(),
onClose:  new _EventDispatch(),

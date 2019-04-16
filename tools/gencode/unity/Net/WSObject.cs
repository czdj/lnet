using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using WebSocketSharp;
using System.Reflection;


[System.Serializable]
class JsonServerData
{
    public int id;
    public string ip;
    public string net;
    public int port;
}

[System.Serializable]
class JsonLoginData
{
    public int error;
    public int roleId;
    public JsonServerData server;
}

public partial class WSObject : MonoBehaviour
{
    private WebSocket ws;
    private string session;
    private int id;
    private string url;
    private bool disable;
    private PropertyInfo[] properInfo;
    public static Proto.ServerTime serverTime = new Proto.ServerTime();
    private List<Proto.S2C> s2cList = new List<Proto.S2C>();
    private object lockObj = new object();
    public bool testLogin;
    public string testLoginName="test1";
    public bool proxy;
    public string authUrl = "http://127.0.0.1:5000";
    public string proxyUrl = "ws://127.0.0.1:5000";
    public IEnumerator CoLogin(string name)
    {
        var www = new WWW(authUrl  + "/fastlogin?name=" + name + "&channel=mine&openid=mine_" + name);
        yield return www;
        var json = JsonUtility.FromJson<JsonLoginData>(www.text);
        if (proxy)
        {
            this.url = proxyUrl + "/proxy";
        }
        else
        {
            this.url = json.server.net + "://" + json.server.ip + ":" + json.server.port.ToString();
        }
        this.id = json.roleId;
        ReConnect();
    }

    private IEnumerator CoHeart()
    {
        yield return new WaitForSeconds(10);
        var c2s = new Proto.C2S();
        var pb = new Proto.ServerTimeC2S();
        pb.id = this.id;
        c2s.serverTimeC2S = pb;
        SendC2S(c2s);
    }

    private void OnOpen(object sender, EventArgs e)
    {
        Debug.Log("websocket connect ok");
        var c2s = new Proto.C2S();
        var pb = new Proto.GamerLoginC2S();
        pb.id = this.id;
        pb.session = this.session;
        c2s.gamerLoginC2S = pb;
        SendC2S(c2s);
    }

    private void OnMessage(object sender, MessageEventArgs e)
    {
        var stream = new System.IO.MemoryStream(e.RawData);
        var s2c = ProtoBuf.Serializer.Deserialize<Proto.S2C>(stream);
        if (s2c.error > 0)
        {
            Debug.LogErrorFormat("recv s2c error:{0}", s2c.error);
        }
        else if (s2c.gamerLoginS2C != null)
        {
            this.session = s2c.gamerLoginS2C.main.session;
        }
        
        lock (this.lockObj)
        {
            this.s2cList.Add(s2c);
        }
    }

    private void OnError(object sender, ErrorEventArgs e)
    {
        Debug.LogError(e.Message);
        this.ReConnect();
    }

    private void OnClose(object sender, CloseEventArgs e)
    {
        Debug.LogError(e.Reason);
        this.ReConnect();
    }

    private void ReConnect()
    {
        if (disable)
        { 
            return;
        }
        if (this.ws != null)
        {
            this.ws.OnOpen -= new EventHandler(OnOpen);
            this.ws.OnMessage -= new EventHandler<MessageEventArgs>(OnMessage);
            this.ws.OnError -= new EventHandler<ErrorEventArgs>(OnError);
            this.ws.OnClose -= new EventHandler<CloseEventArgs>(OnClose);
        }
        this.ws = new WebSocket(this.url);
        this.ws.OnOpen += new EventHandler(OnOpen);
        this.ws.OnMessage += new EventHandler<MessageEventArgs>(OnMessage);
        this.ws.OnError += new EventHandler<ErrorEventArgs>(OnError);
        this.ws.OnClose += new EventHandler<CloseEventArgs>(OnClose);
        this.ws.Connect();
    }


    public void SendC2S(Proto.C2S c2s)
    {
        var stream = new System.IO.MemoryStream();
        ProtoBuf.Serializer.Serialize(stream, c2s);
        ws.Send(stream.ToArray());
    }

    private void OnDisable()
    {
        disable = true;
    }

    void Start()
    {
        var type = typeof(Proto.S2C);
        properInfo = type.GetProperties();
        DontDestroyOnLoad(gameObject);
        if (testLogin)
        {
            StartCoroutine(CoLogin(testLoginName));
        }
    }

    void Update()
    {
        while (this.s2cList.Count > 0)
        {
            Proto.S2C s2c = null;
            lock (this.lockObj)
            {
                s2c = this.s2cList[0];
                this.s2cList.RemoveAt(0);
            }
            if (s2c == null)
            {
                continue;
            }
            if (s2c.error == 0)
            {
                if (s2c.gamerLoginS2C != null)
                {
                    this.session = s2c.gamerLoginS2C.main.session;
                    Debug.LogFormat("login to logic ok reconn session:{0}", this.session);
                    StopAllCoroutines();
                    StartCoroutine(CoHeart());
                }
                else if (s2c.serverTimeS2C != null)
                {
                    WSObject.serverTime = s2c.serverTimeS2C.time;
                    StartCoroutine(CoHeart());
                }
            }

            foreach (var info in properInfo)
            {
                if (info.Name == "error" || info.Name == "key")
                {
                    continue;
                }
                if (info.GetValue(s2c, null) != null)
                {
                    Debug.LogError(info.Name);
                }
            }
        }
    }
}


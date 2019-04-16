using UnityEngine;
using System;
using System.IO;
using System.Net;
using System.Net.Sockets;
using System.Collections.Generic;
using System.Net.NetworkInformation;
using System.Runtime.InteropServices;
using System.Threading;



public class UdpData
{
    public byte[] data;
    public UdpData()
    {
        data = new byte[1 << 12];
    }
}

public class UdpObject : MonoBehaviour
{
#if UNITY_ANDROID || UNITY_STANDALONE
    [DllImport("udp")]
    private static extern int udp_connect(string ip, int port);

    [DllImport("udp")]
    private static extern int udp_send(ref byte buffer, int len);

    [DllImport("udp")]
    private static extern int udp_recv(ref byte buffer, int len);

    [DllImport("udp")]
    private static extern bool udp_available();

    [DllImport("udp")]
    private static extern bool udp_continue();

    [DllImport("udp")]
    private static extern void udp_close();
#elif UNITY_IOS

    [DllImport("__Internal")]
    private static extern int udp_connect(string ip, int port);

    [DllImport("__Internal")]
    private static extern int udp_send(ref byte buffer, int len);

    [DllImport("__Internal")]
    private static extern int udp_recv(ref byte buffer, int len);

    [DllImport("__Internal")]
    private static extern bool udp_available();

    [DllImport("__Internal")]
    private static extern bool udp_continue();

    [DllImport("__Internal")]
    private static extern void udp_close();
#endif
    public static UdpObject instance;
    public static int reconnTimeout;

    protected Socket udpClient;
    protected EndPoint serverAddr;
       

    private UdpData udpData = new UdpData();
    private bool _threadRecv = false;
    private string _ip;
    private uint _port;
    private bool ipv6;
    private Thread _thread;
    private int _lastRecv = 0;
    private int _tick = 0;
#if (UNITY_ANDROID || UNITY_IOS) && !UNITY_EDITOR
    public const bool useDLL = true;
#else
    public const bool useDLL = false;
#endif
    public System.Action<UdpData> threadRecv;
    private bool _threadStop = false;
    void Awake()
    {
        instance = this;
    }

    /// <summary>
    /// 使用SDK加速网络  重连   在主线程调用
    /// </summary>
    /// <param name="reConnect">是否是重连</param>
    private void StartAddSpeed(bool reConnect)
    {
        ReConnect();
    }

    private void ReConnect()
    {
        try
        {
            if (useDLL)
            {
                udp_connect(_ip, (int) _port);
            }
            else
            {
                if (udpClient != null)
                {
                    udpClient.Close();
                }
                udpClient = new Socket(ipv6 ? AddressFamily.InterNetworkV6 : AddressFamily.InterNetwork, SocketType.Dgram, ProtocolType.Udp);
                udpClient.SendBufferSize = 20480;
                udpClient.ReceiveBufferSize = 20480;
                udpClient.Blocking = false;
                udpClient.Connect(serverAddr);
            }
        }
        catch (Exception e)
        {
            Debug.LogErrorFormat("ReConnect  exception! {0},{1}", e.ToString(), e.StackTrace);
        }
    }

    public void Connect(string ip, uint port, bool threadRecv, int reconnTimeout)
    {
        _tick = 0;
        _port = port;
        _threadRecv = threadRecv;
        IPV6Help.GetIPType(ip, (int)port, out _ip, out ipv6);
        Debug.LogFormat("connect to pvp server {0}:{1} thread recv:{2}", ip, port, threadRecv);
        serverAddr = new IPEndPoint(IPAddress.Parse(_ip), (int)_port);
        Stop(false);

        StartAddSpeed(false);
        if (_threadRecv)
        {
            _thread = new Thread(new ThreadStart(ThreadRecv));
            _thread.IsBackground = true;
            _thread.Start();
        }
    }

    public void Stop(bool destroy = true)
    {
        if (_thread != null)
        {
            _threadStop = true;
            _thread.Join();
            _thread = null;
            _threadStop = false;
        }
        if (useDLL)
        {
            udp_close();
        }
        else
        {
            if (udpClient != null)
            {
                udpClient.Close();
                udpClient = null;
            }
        }

        if (destroy)
        {
            if (instance != null)
            {
                instance = null;
                Destroy(gameObject);
            }
        }
    }

    public void OnDisable()
    {
        Stop(false);
    }

    public bool ThreadRecvData(Action<UdpData> onRecv)
    {
        int len = 0;
        int prevLen = 0;
        try
        {
            while ((len = useDLL ? udp_recv(ref udpData.data[0], udpData.data.Length) :
                udpClient.Receive(udpData.data, udpData.data.Length, SocketFlags.None)) > 0)
            {
                prevLen = len;
            }
            len = prevLen;
        }
        catch (Exception)
        {
        }
        if (len > 0)
        {
            onRecv(udpData);
            _lastRecv = _tick;
            return true;
        }
        return false;
    }

    public void ThreadRecv()
    {
        try
        {
            while (!_threadStop)
            {
                Thread.Sleep(1);
                ThreadRecvData(threadRecv);
                threadRecv(null);
                _tick++;
                if (_tick - _lastRecv > UdpObject.reconnTimeout)
                {
                    _lastRecv = _tick;
                    ReConnect();
                    //Loom.QueueOnMainThread((param) => { StartAddSpeed(true); }, null);
                }
            }
        }
        catch (Exception e)
        {
            Debug.LogErrorFormat("Logic Loop Failed {0} {1} \n{3}", e.Message, instance, e.StackTrace);
        }
    }

    public void GetRecv(Action<UdpData> onRecv)
    {
        while (udpClient.Available > 0)
        {
            int len = 0;
            if (useDLL)
            {
                len = udp_recv(ref udpData.data[0], udpData.data.Length);
            }
            else
            {
                len = udpClient.Receive(udpData.data, udpData.data.Length, SocketFlags.None);
            }

            if (len > 0)
            {
                onRecv(udpData);
            }
        }
    }

    public void Send(byte[] data, int len)
    {
        try
        {
            if (useDLL)
            {
                udp_send(ref data[0], len);
            }
            else
            {
                udpClient.Send(data, len, SocketFlags.None);
            }
        }
        catch (Exception e)
        {
        }
    }
    public bool Available
    {
        get
        {
            return useDLL ? udp_available() : udpClient != null;
        }
    }

}

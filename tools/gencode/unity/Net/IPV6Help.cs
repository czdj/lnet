using System;
using System.Collections;
using System.Collections.Generic;
using System.Runtime.InteropServices;
using UnityEngine;

public class IPV6Help
{

#if UNITY_IPHONE && !UNITY_EDITOR
    [DllImport("__Internal")]
    private static extern bool supportIPV6();

    [DllImport("__Internal")]
    private static extern string getIPv6(string mHost, string mPort);
#endif

    //"192.168.1.1&&ipv4"
    public static string GetIPv6(string mHost, string mPort)
    {
#if UNITY_IPHONE && !UNITY_EDITOR
			if (supportIPV6() == true)
			{
				string mIPv6=getIPv6(mHost, mPort);
				return mIPv6;
			}
#endif

        return mHost + "&&ipv4";

    }

    public static void GetIPType(string serverIp, int serverPorts, out string newServerIp, out bool ipv6)
    {
        ipv6 = false;
        newServerIp = serverIp;

        try
        {
            string mIPv6 = GetIPv6(serverIp, serverPorts.ToString());
            if (!string.IsNullOrEmpty(mIPv6))
            {
                string[] m_StrTemp = System.Text.RegularExpressions.Regex.Split(mIPv6, "&&");
                if (m_StrTemp != null && m_StrTemp.Length >= 2)
                {
                    string IPType = m_StrTemp[1];
                    if (IPType == "ipv6")
                    {
                        newServerIp = m_StrTemp[0];
                        ipv6 = true;
                    }
                }
            }
        }
        catch (Exception e)
        {
            Debug.LogError("GetIPv6 error:" + e);
        }
    }
}
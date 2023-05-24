const prxSelectedProxy = "${PRX_PROXY}";
const prxMListener = function() {
    console.log("Proxy extension start. Selected proxy: ", prxSelectedProxy);
    let proxy = prxSelectedProxy.split("@");
    let proxyAddr = "";
    if(proxy.length > 1) {
        proxyAddr = proxy[1].split(":");
    } else {
        proxyAddr = proxy[0].split(":");
    }
    
    var config = {
        mode: "fixed_servers",
        rules: {
            singleProxy: {
                scheme: "http",
                host: proxyAddr[0],
                port: parseInt(proxyAddr[1])
            },
            bypassList: ["foobar.com"]
        }
    };

    chrome.proxy.settings.set(
        {value: config, scope: "regular"}, 
        function() {}
    );

    console.log("Proxy extension initialized with conf", config);

};

function callbackFn(details) {
    console.log("ProxyCallback called details:", details);
    let proxy = prxSelectedProxy.split("@");
    let proxyAuth = proxy[0].split(":");

    return {
        authCredentials: {
            username: proxyAuth[0],
            password: proxyAuth[1]
        }
    };
}

// Setup
chrome.runtime.onStartup.addListener(prxMListener);
chrome.runtime.onInstalled.addListener(prxMListener);
chrome.webRequest.onAuthRequired.addListener(
    callbackFn,
    {urls: ["<all_urls>"]},
    ['blocking']
);
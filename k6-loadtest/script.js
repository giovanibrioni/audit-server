import http from 'k6/http';
import { check, sleep } from 'k6';
export let options = {
    vus: 1,
    duration: '60s',
};


export default function () {
    var params = {
        headers: {
          'Content-Type': 'application/json',
        },
    };

    const payload = JSON.stringify([{"route":{"ws_id":"899e0d91-2b05-495e-bf74-8426a6b43155","preserve_host":false,"paths":["/book-store"],"regex_priority":0,"request_buffering":true,"response_buffering":true,"https_redirect_status_code":426,"name":"book-store-route","id":"a510ebe5-ba62-4c94-a0b1-d020fa07a57a","protocols":["http","https"],"path_handling":"v1","service":{"id":"4bce5799-3127-492f-baae-e2672c0a1e55"},"created_at":1666077712,"updated_at":1666077712,"strip_path":true},"started_at":1666449742317,"client_ip":"10.244.1.0","tries":[{"balancer_latency":0,"ip":"129.153.98.232","port":8080,"balancer_start":1666449742358}],"service":{"ws_id":"899e0d91-2b05-495e-bf74-8426a6b43155","connect_timeout":60000,"protocol":"http","read_timeout":60000,"enabled":true,"host":"129.153.98.232","write_timeout":60000,"name":"book-store-service","retries":5,"id":"4bce5799-3127-492f-baae-e2672c0a1e55","tags":[],"port":8080,"created_at":1666077688,"updated_at":1666077920},"response":{"status":200,"size":417,"headers":{"x-kong-upstream-latency":"3","date":"Sat, 22 Oct 2022 14:42:22 GMT","content-length":"198","content-type":"application/json; charset=utf-8","x-kong-proxy-latency":"41","via":"kong/2.8.1","connection":"close"},"body":{"data":[{"author":"J. K. Rowling","id":"1","title":"Harry Potter AAA"},{"author":"J. R. R. Tolkien","id":"2","title":"The Lord of the Rings"},{"author":"L. Frank Baum","id":"3","title":"The Wizard of Oz"}],"type":"json"}},"latencies":{"request":44,"proxy":3,"kong":41},"request":{"querystring":{},"size":155,"method":"GET","uri":"/book-store/books","headers":{"connection":"keep-alive","accept-encoding":"gzip, deflate","host":"api-px.pocs.dev.br","accept":"*/*","user-agent":"HTTPie/3.2.1"},"url":"http://api-px.pocs.dev.br:80/book-store/books"},"upstream_uri":"/books"}]);


    let res = http.post('http://host.docker.internal:8080/audit/batch', payload, params)
    //console.log(`Status: ${res.status}, Body: ${res.body}`);
    check(res, { 'status was 200': (r) => r.status == 200 });
    //sleep(1);
}

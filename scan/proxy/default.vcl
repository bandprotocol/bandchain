vcl 4.0;

backend default {
  .host = "172.18.0.20";
  .port = "1317";
}

backend bandsv {
  .host = "172.18.0.16";
  .port = "5000";
}

sub vcl_deliver {
  if (req.url ~ "/") {
    set resp.http.Access-Control-Allow-Origin = "*";
    set resp.http.Access-Control-Allow-Methods = "GET, OPTIONS, POST, PATCH, PUT, DELETE";
    set resp.http.Access-Control-Allow-Headers = "Origin, Accept, Content-Type, X-Requested-With, X-CSRF-Token";
  }
}

sub vcl_recv {
  if (req.url ~ "^/bandsv/") {
    set req.url = regsub(req.url, "^/bandsv", "/");
    set req.backend_hint = bandsv;
  } else {
    set req.backend_hint = default;
  }
}

sub vcl_backend_response {
  set beresp.ttl = 3s;
}

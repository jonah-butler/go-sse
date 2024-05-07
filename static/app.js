const es = new EventSource("http://localhost:8080/event/sse?userId=771996");

es.onopen = (...args) => {
  console.log("event source opened...", args);
};

es.onmessage = (msg) => {
  console.log("received message: ", msg);
};

es.onerror = (err) => {
  console.log("got error: ", err);
  es.close();
};

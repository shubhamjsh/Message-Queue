import http from "k6/http";
import { sleep } from "k6";
import exec from "k6/execution"

export const options = {
    duration: "1m",
    vus: 50,
    thresholds: {
        http_req_failed: ["rate<0.01"], // http errors should be less than 1%
        http_req_duration: ["p(95)<500"], // 95 percent of response times must be below 500ms
    },
};

const endpoints = [
    "http://localhost:4001/webhook",
]
const topicSets = [
    ["orders"], ["payments"], ["metadata"], ["notification"]
]

export default function () {
    const epIndex = (__VU - 1) % endpoints.length;
    const tIndex = (__VU - 1) % topicSets.length
    const res = http.post("http://localhost:3000/v1/subscribe",
        JSON.stringify({
            endpoint: endpoints[epIndex],
            topics: topicSets[tIndex],
        }),
        { headers: { "Content-Type": "application/json" }});

    abortOnFailure("Subscribe", res)
    sleep(1);

    const pubRes = http.post("http://localhost:3000/v1/publish",
        JSON.stringify({
            topic: "orders",
            message: "order123 published"
        }),
        { headers: { "Content-Type": "application/json"}})

    abortOnFailure("Publish", pubRes)
}

function abortOnFailure(label, res) {
    if (res.status < 200 || res.status >= 300) {
        exec.test.abort(
            `${label} failed: status=${res.status}, body=${res.body}, error=${res.error}`        );
    }
}

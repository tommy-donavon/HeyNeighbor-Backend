import consul from "consul";
import os from "os";

const client = consul({ host: "consul", port: "8500" });


export function newConsulDetails(
  serviceName,
  port
) {
  return {
    name: serviceName,
    id: os.hostname(),
    address: os.hostname(),
    port: Number(port),
    check: {
      ttl: '10s',
      deregister_critical_service_after: '1m'
    },
  };
}

export function registerService(details) {
  client.agent.service.register(details, (err) => {
    if (err) {
      throw new Error(err.message);
    }

    setInterval(() => {
      client.agent.check.pass({id:`service:${details.id}`}, err => {
        if (err) throw new Error(err);
      });
    }, 5 * 1000);

  });
}

export function deregisterService(id) {
  client.agent.service.deregister({id: id}, (err) => {
    if (err) console.error(err.message);
    console.log(`de-registered`);
  });
}

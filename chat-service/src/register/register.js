import consul from 'consul';
import os from 'os';
import fetch from 'node-fetch';
import Resilient from 'resilient';

const client = consul({ host: 'consul', port: '8500' });

export function newConsulDetails(serviceName, port) {
  return {
    name: serviceName,
    id: os.hostname(),
    address: os.hostname(),
    port: Number(port),
    check: {
      ttl: '10s',
      deregister_critical_service_after: '1m',
    },
  };
}

export function registerService(details) {
  client.agent.service.register(details, (err) => {
    if (err) {
      throw new Error(err.message);
    }

    setInterval(() => {
      client.agent.check.pass({ id: `service:${details.id}` }, (err) => {
        if (err) throw new Error(err);
      });
    }, 5 * 1000);
  });
}

export function deregisterService(id) {
  client.agent.service.deregister({ id: id }, (err) => {
    if (err) console.error(err.message);
    console.log(`de-registered`);
  });
}

export async function lookupService(serviceName) {
  try {
    const request = await fetch(`http://consul:8500/v1/agent/services`);
    const result = await request.json();
    for (const service in result) {
      if (serviceName === result[service]['Service']) {
        return `http://${result[service]['Address']}:${result[service]['Port']}/`;
      }
    }
  } catch (err) {
    console.error(err);
  }
  return "";
}

const endpoints = [
  process.env.ARM_FLOAT_GO,
  process.env.ARM_INT_GO,
  process.env.INTEL_FLOAT_GO,
  process.env.INTEL_INT_GO,
  process.env.ARM_FLOAT_JS,
  process.env.ARM_INT_JS,
  process.env.INTEL_FLOAT_JS,
  process.env.INTEL_INT_JS,
];

const burstDuration = 15 * 60 * 1000; // 15 minutes in Milliseconds = 15m * 60s * 1000ms

const makeRequest = async (endpoint) => {
  // await fetch endpoint
};

const currentTime = new Date();

const burstEndpoints = async (startTime) => {
  while (new Date() - startTime < burstDuration) {
    const res = await Promise.all(endpoints.map(makeRequest));
    const data = res.map((res) => res.data);
  }
};

burstEndpoints(currentTime);

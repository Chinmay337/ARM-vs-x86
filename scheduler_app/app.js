const ARM_FLOAT_GO = process.env.ARM_FLOAT_GO;
const ARM_INT_GO = process.env.ARM_INT_GO;
const INTEL_FLOAT_GO = process.env.INTEL_FLOAT_GO;
const INTEL_INT_GO = process.env.INTEL_INT_GO;
const ARM_FLOAT_JS = process.env.ARM_FLOAT_JS;
const ARM_INT_JS = process.env.ARM_INT_JS;
const INTEL_FLOAT_JS = process.env.INTEL_FLOAT_JS;
const INTEL_INT_JS = process.env.INTEL_INT_JS;

const endpoints = [
  ARM_FLOAT_GO,
  ARM_INT_GO,
  INTEL_FLOAT_GO,
  INTEL_INT_GO,
  ARM_FLOAT_JS,
  ARM_INT_JS,
  INTEL_FLOAT_JS,
  INTEL_INT_JS,
];

const makeRequest = async (endpoint) => {
  // await fetch endpoint
};

const currentTime = new Date();

const burstEndpoints = async (startTime) => {
  while (new Date() - startTime < 15) {
    const res = await Promise.all(endpoints.map(makeRequest));
    const data = res.map((res) => res.data);
  }
};

burstEndpoints(currentTime);

const ComputeInt = (input) => {
  //sample_arr := []int{10, 11, 192, 200, 12}

  input.forEach((val) => {
    let counter = 0;
    for (let i = 2; i <= val; i++) {
      const ok = isPrime(i);
      if (ok) {
        counter++;
      }
    }
    //  console.log("There are ", counter, "prime numbers upto ", val);
  });
};

const isPrime = (num) => {
  let i = 2;

  while (i <= num) {
    let rem = num % i;
    if (rem === 0) {
      return false;
    }
    if (i * i > num) {
      return true;
    }
    i += 1;
  }
  return true;
};

const prepareReturnValue = (code, content) => {
  return {
    statusCode: code,
    body: content,
  };
};

//-------------------------------------------------------------
//  Main Lambda handler
//-------------------------------------------------------------

exports.handler = async (event) => {
  console.log("request: " + JSON.stringify(event));

  if (event.body) {
    const ele = JSON.parse(event.body);

    if (Array.isArray(ele?.input)) {
      const start = new Date();

      ComputeInt(ele.input);

      stop = new Date();

      var mstime = stop.getTime() - start.getTime();
      const msg = "Elapsed time in seconds: " + (mstime / 1000).toFixed(3);

      return prepareReturnValue(200, msg);
    }
  }

  return prepareReturnValue(400, "Did not run benchmarks");
};

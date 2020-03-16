let findNumberOfPage = (amount, limit) =>
  if (amount mod limit == 0) {
    amount / limit;
  } else {
    amount / limit + 1;
  };

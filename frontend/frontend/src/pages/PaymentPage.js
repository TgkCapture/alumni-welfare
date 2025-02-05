import React from "react";
import { Button, TextField, Typography } from "@mui/material";

function PaymentPage() {
  return (
    <div>
      <Typography variant="h5">Make a Payment</Typography>
      <TextField label="Full Name" variant="outlined" fullWidth sx={{ mt: 2 }} />
      <TextField label="Amount" type="number" variant="outlined" fullWidth sx={{ mt: 2 }} />
      <Button variant="contained" color="primary" sx={{ mt: 2 }}>
        Pay Now
      </Button>
    </div>
  );
}

export default PaymentPage;

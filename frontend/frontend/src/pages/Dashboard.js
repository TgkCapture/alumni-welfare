import React from "react";
import { Button, Typography } from "@mui/material";
import { useNavigate } from "react-router-dom";

function Dashboard() {
  const navigate = useNavigate();

  return (
    <div>
      <Typography variant="h5">Welcome to the Alumni Welfare System</Typography>
      <Button variant="contained" color="primary" onClick={() => navigate("/payment")} sx={{ mt: 2 }}>
        Make a Payment
      </Button>
    </div>
  );
}

export default Dashboard;

import React from "react";
import { Button, Typography, Container, Box } from "@mui/material";
import { useNavigate } from "react-router-dom";

function Dashboard() {
  const navigate = useNavigate();

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 5, textAlign: "center" }}>
        <Typography variant="h5">Welcome to the Alumni Welfare System</Typography>
        <Button
          variant="contained"
          color="primary"
          onClick={() => navigate("/payment")}
          sx={{ mt: 2 }}
        >
          Make a Payment
        </Button>
      </Box>
    </Container>
  );
}

export default Dashboard;

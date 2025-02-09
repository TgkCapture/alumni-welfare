import React, { useState } from "react";
import { Button, TextField, Typography, Container, Box } from "@mui/material";
import { makePayment } from "../services/api";

function PaymentPage() {
  const [name, setName] = useState("");
  const [amount, setAmount] = useState("");
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");

  const handlePayment = async () => {
    if (!name || !amount) {
      alert("Please enter your name and amount.");
      return;
    }

    setLoading(true);
    setMessage("");

    try {
      const response = await makePayment(name, amount);
      setMessage(response.message || "Payment successful!");
    } catch (error) {
      setMessage("Payment failed. Please try again.");
    }

    setLoading(false);
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 5, textAlign: "center" }}>
        <Typography variant="h5" gutterBottom>
          Make a Payment
        </Typography>
        <TextField
          fullWidth
          label="Full Name"
          variant="outlined"
          value={name}
          onChange={(e) => setName(e.target.value)}
          sx={{ my: 2 }}
        />
        <TextField
          fullWidth
          label="Amount"
          variant="outlined"
          type="number"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
          sx={{ my: 2 }}
        />
        <Button
          variant="contained"
          color="primary"
          onClick={handlePayment}
          disabled={loading}
        >
          {loading ? "Processing..." : "Pay Now"}
        </Button>
        {message && (
          <Typography variant="body1" sx={{ mt: 2, color: "green" }}>
            {message}
          </Typography>
        )}
      </Box>
    </Container>
  );
}

export default PaymentPage;

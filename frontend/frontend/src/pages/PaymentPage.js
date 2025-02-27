import React, { useState, useEffect } from "react";
import { Button, TextField, Typography, Container, CircularProgress, Alert } from "@mui/material";
import axios from "axios";

function PaymentPage() {
  const [name, setName] = useState("");
  const [amount, setAmount] = useState("");
  const [month, setMonth] = useState("");
  const [history, setHistory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [paymentSuccess, setPaymentSuccess] = useState(null);

  const token = localStorage.getItem("token");

  useEffect(() => {
    axios
      .get(`${API_BASE_URL}/payments  `, { headers: { Authorization: token } })
      .then((res) => setHistory(res.data))
      .catch((err) => setError("Failed to load payment history."))
      .finally(() => setLoading(false));
  });

  const handlePayment = async () => {
    if (!name || !amount || !month) {
      setError("All fields are required.");
      return;
    }

    const paymentData = { name, amount: parseInt(amount), month: parseInt(month) };

    if (paymentData.amount <= 0) {
      setError("Amount must be greater than zero.");
      return;
    }

    if (paymentData.month < 1 || paymentData.month > 12) {
      setError("Month must be between 1 and 12.");
      return;
    }

    try {
      await axios.post("http://localhost:8080/pay", paymentData, { headers: { Authorization: token } });
      setPaymentSuccess("Payment successful!");
      setError(null);
      
      // Refresh payment history
      const res = await axios.get("http://localhost:8080/payments", { headers: { Authorization: token } });
      setHistory(res.data);
      
    } catch (err) {
      setError("Payment failed. Please try again.");
    }
  };

  return (
    <Container>
      <Typography variant="h5">Make a Payment</Typography>

      {error && <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>}
      {paymentSuccess && <Alert severity="success" sx={{ mt: 2 }}>{paymentSuccess}</Alert>}

      <TextField
        label="Full Name"
        variant="outlined"
        fullWidth
        sx={{ mt: 2 }}
        value={name}
        onChange={(e) => setName(e.target.value)}
      />
      <TextField
        label="Amount"
        type="number"
        variant="outlined"
        fullWidth
        sx={{ mt: 2 }}
        value={amount}
        onChange={(e) => setAmount(e.target.value)}
      />
      <TextField
        label="Month (1-12)"
        type="number"
        variant="outlined"
        fullWidth
        sx={{ mt: 2 }}
        value={month}
        onChange={(e) => setMonth(e.target.value)}
      />
      <Button variant="contained" color="primary" fullWidth sx={{ mt: 2 }} onClick={handlePayment}>
        Pay Now
      </Button>

      <Typography variant="h6" sx={{ mt: 4 }}>Payment History</Typography>
      {loading ? <CircularProgress sx={{ mt: 2 }} /> : 
        history.length > 0 ? history.map((pay, index) => (
          <Typography key={index} sx={{ mt: 1 }}>
            <strong>{pay.name}</strong> paid <strong>${pay.amount}</strong> for month <strong>{pay.month}</strong>
          </Typography>
        )) : <Typography sx={{ mt: 2 }}>No payment history available.</Typography>
      }
    </Container>
  );
}

export default PaymentPage;

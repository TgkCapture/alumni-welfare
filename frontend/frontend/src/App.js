import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { Container, CssBaseline, Typography } from "@mui/material";
import PaymentPage from "./pages/PaymentPage";
import Dashboard from "./pages/Dashboard";

function App() {
  return (
    <Router>
      <CssBaseline />
      <Container>
        <Typography variant="h4" gutterBottom>Alumni Welfare Payment System</Typography>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/payment" element={<PaymentPage />} />
        </Routes>
      </Container>
    </Router>
  );
}

export default App;

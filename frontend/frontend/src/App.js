import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { Container, CssBaseline, Typography } from "@mui/material";
import ProtectedRoute from "./components/ProtectedRoute";
import PaymentPage from "./pages/PaymentPage";
import Dashboard from "./pages/Dashboard";
import Register from "./pages/Register";
import Login from "./pages/Login";

function App() {
  return (
    <Router>
      <CssBaseline />
      <Container>
        <Typography variant="h4" gutterBottom>
          Alumni Welfare Payment System
        </Typography>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/register" element={<Register />} />
          <Route path="/login" element={<Login />} />

          {/* Protect these routes */}
          <Route element={<ProtectedRoute />}>
            <Route path="/payment" element={<PaymentPage />} />
          </Route>
        </Routes>
      </Container>
    </Router>
  );
}

export default App;

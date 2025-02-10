import React, { useState } from "react";
import { TextField, Button, Typography, Container } from "@mui/material";
import { useNavigate } from "react-router-dom";

function Register() {
  const [formData, setFormData] = useState({ name: "", email: "", password: "" });
  const navigate = useNavigate();

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const response = await fetch("http://localhost:8080/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(formData),
    });

    if (response.ok) {
      navigate("/login");
    } else {
      console.error("Registration failed");
    }
  };

  return (
    <Container>
      <Typography variant="h5">Register</Typography>
      <form onSubmit={handleSubmit}>
        <TextField label="Full Name" name="name" onChange={handleChange} fullWidth required sx={{ mt: 2 }} />
        <TextField label="Email" name="email" type="email" onChange={handleChange} fullWidth required sx={{ mt: 2 }} />
        <TextField label="Password" name="password" type="password" onChange={handleChange} fullWidth required sx={{ mt: 2 }} />
        <Button type="submit" variant="contained" color="primary" fullWidth sx={{ mt: 2 }}>Register</Button>
      </form>
    </Container>
  );
}

export default Register;

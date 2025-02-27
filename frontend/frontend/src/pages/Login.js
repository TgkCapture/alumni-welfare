import React, { useState } from "react";
import { TextField, Button, Typography, Container } from "@mui/material";
import { useNavigate } from "react-router-dom";

const API_BASE_URL = process.env.REACT_APP_API_BASE_URL;

function Login() {
  const [formData, setFormData] = useState({ email: "", password: "" });
  const navigate = useNavigate();

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const response = await fetch(`${API_BASE_URL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
    });

    if (response.ok) {
        const data = await response.json();
        localStorage.setItem("token", data.token);
        navigate("/");
    } else {
        console.error("Login failed");
    }
};

  return (
    <Container>
      <Typography variant="h5">Login</Typography>
      <form onSubmit={handleSubmit}>
        <TextField label="Email" name="email" type="email" onChange={handleChange} fullWidth required sx={{ mt: 2 }} />
        <TextField label="Password" name="password" type="password" onChange={handleChange} fullWidth required sx={{ mt: 2 }} />
        <Button type="submit" variant="contained" color="primary" fullWidth sx={{ mt: 2 }}>Login</Button>
      </form>
    </Container>
  );
}

export default Login;

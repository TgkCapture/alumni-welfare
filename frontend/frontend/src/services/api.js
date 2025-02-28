const API_BASE_URL = process.env.REACT_APP_API_BASE_URL;

export const makePayment = async (name, amount) => {
  try {
    const response = await fetch(`${API_BASE_URL}/pay`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name, amount: parseInt(amount, 10) }),
    });

    return await response.json();
  } catch (error) {
    console.error("Payment failed:", error);
    throw error;
  }
};

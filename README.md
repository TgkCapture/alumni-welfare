# Project Report: Alumni Class Welfare Payment System

## 1. Introduction
The **Alumni Class Welfare Payment System** is a web-based application designed to streamline the process of collecting monthly welfare contributions from alumni members. The system automates payment tracking, generates proof of payment, and updates records in real-time, enhancing transparency and efficiency.

## 2. Project Objectives
- Automate the collection of monthly welfare contributions.  
- Provide an intuitive web interface for members to register, log in, and make payments.  
- Generate automatic proof of payment receipts.  
- Update an Excel sheet or database with payment records in real-time.  
- Offer WhatsApp notifications for payment confirmations (future scope).  

## 3. System Overview
The system is built using the following technologies:  
- Backend: **Go (Golang)** with **Gin framework**  
- Frontend: **React.js** with **MUI (Material UI)** 
- Database: **PostgreSQL** with **GORM ORM**  
- Authentication: JWT-based user authentication  
- Payment Gateway: Integrated API (Paychangu)  

### Features
- User registration and login  
- Payment management interface  
- Month-based payment selection  
- Payment confirmation and receipt generation  
- Payment history view  
- Admin interface (future scope)  

## 4. Implementation Approach
### 4.1 Backend Development
The backend API handles user authentication, payment processing, and data storage. Key endpoints include:  
- `/register`: User registration  
- `/login`: User login  
- `/pay`: Payment processing  
- `/payments`: Retrieve payment history  

### 4.2 Frontend Development
The frontend provides an intuitive interface where users can:
- Register and log in  
- View their payment status  
- Select the month for payment  
- Initiate and confirm payments  

### 4.3 Payment Gateway Integration
The system integrates with a payment gateway to enable secure online payments. Each transaction generates a unique reference and proof of payment.

## 5. Future Enhancements
- WhatsApp notifications for payment confirmations  
- Admin dashboard for payment oversight  
- Bulk payment uploads  
- Multi-user access levels  

## 6. Conclusion
The **Alumni Class Welfare Payment System** is a step towards digitizing alumni welfare contributions, improving efficiency, and enhancing member experience. The project lays a solid foundation for future scalability and additional features.


---
**Prepared by:** Tawonga Grant Kanyenda  
**Date:** 27 February 2025


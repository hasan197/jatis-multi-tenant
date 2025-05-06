import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { RabbitMQPublisher } from './pages/RabbitMQPublisher';

const App: React.FC = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<RabbitMQPublisher />} />
        <Route path="/dashboard" element={<RabbitMQPublisher />} />
      </Routes>
    </Router>
  );
};

export default App; 
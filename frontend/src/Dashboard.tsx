import { useState, useEffect } from "react";
import { Header } from './component/Header';
import { Summary } from './component/Summary';
import { ExpensesTable } from './component/Expense';
import { CreateExpense } from './component/CreateExpense';

interface Expense {
  id: number;
  title: string;
  description: string;
  category: string;
  amount: number;
  date: Date;
}

export const Dashboard = () => {
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [categories, setCategories] = useState<string[]>(["All"]);

  const fetchData = async () => {
    try {
      const res = await fetch("http://localhost:5000/api/v1/expenses", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: localStorage.getItem("token") || "",
        },
      });
      if (!res.ok) throw new Error("Failed to fetch expenses");

      const response = await res.json();
      const expensesData: Expense[] = response.data?.data || [];

      setExpenses(expensesData);

      const uniqueCategories = Array.from(
        new Set(expensesData.map((e) => e.category))
      );
      setCategories(["All", ...uniqueCategories]);
    } catch (error) {
      console.error("Error fetching expenses:", error);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  return (
    <>
      <Header />
      <div style={{ padding: '2rem', background: '#f9fafb', minHeight: '100vh' }}>
        {/* <Summary pieData={pieData} barData={barData} /> */}
        <ExpensesTable expenses={expenses} fetchExpenses={fetchData} categories={categories} />
        <CreateExpense onExpenseCreated={fetchData} />
      </div>
    </>
  );
};

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

interface CategoryDatum {
  name: string;
  value: number;
}
interface TimeDatum {
  name: string;
  total: number;
}

export const Dashboard = () => {
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [categories, setCategories] = useState<string[]>(["All"]);
  const [pieData, setPieData] = useState<CategoryDatum[]>([]);
  const [barDataMonth, setBarDataMonth] = useState<TimeDatum[]>([]);
  const [barDataWeek, setBarDataWeek] = useState<TimeDatum[]>([]);
  const headers = {
    "Content-Type": "application/json",
    Authorization: localStorage.getItem("token") || "",
  };

  const fetchExpenses = async () => {
    try {
      const res = await fetch("http://localhost:5000/api/v1/expenses", { headers });
      if (!res.ok) throw new Error("Failed to fetch expenses");

      const response = await res.json();
      const expensesData: Expense[] = response.data?.data || [];

      setExpenses(expensesData);

      const uniqueCategories = Array.from(new Set(expensesData.map((e) => e.category)));
      setCategories(["All", ...uniqueCategories]);

      fetchCategoryTotals();
      fetchMonthlyTotals();
      fetchWeeklyTotals();
    } catch (error) {
      console.error("Error fetching expenses:", error);
    }
  };

  const fetchCategoryTotals = async () => {
    try {
      const res = await fetch("http://localhost:5000/api/v1/expenses/categories", { headers });
      const response = await res.json();
      const formatted: CategoryDatum[] = response.data || [];
      setPieData(formatted);
    } catch (error) {
      console.error("Error fetching category totals:", error);
    }
  };

  const fetchMonthlyTotals = async () => {
    try {
      const res = await fetch("http://localhost:5000/api/v1/expenses/month", { headers });
      const response = await res.json();
      const formatted: TimeDatum[] = response.data || [];
      setBarDataMonth(formatted);
    } catch (error) {
      console.error("Error fetching monthly totals:", error);
    }
  };

  const fetchWeeklyTotals = async () => {
    try {
      const res = await fetch("http://localhost:5000/api/v1/expenses/week", { headers });
      const response = await res.json();
      const formatted: TimeDatum[] = response.data || [];
      setBarDataWeek(formatted);
    } catch (error) {
      console.error("Error fetching weekly totals:", error);
    }
  };

  useEffect(() => {
    fetchExpenses();
    fetchCategoryTotals();
    fetchMonthlyTotals();
    fetchWeeklyTotals();
  }, []);

  return (
    <>
      <Header />
      <div style={{ padding: '2rem', background: '#1e1e2f', minHeight: '100vh' }}>
        <Summary pieData={pieData} barDataMonth={barDataMonth} barDataWeek={barDataWeek} />
        <ExpensesTable expenses={expenses} fetchExpenses={fetchExpenses} categories={categories} />
        <CreateExpense onExpenseCreated={fetchExpenses} />
      </div>
    </>
  );
};

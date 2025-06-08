import React from "react";
import {
  PieChart,
  Pie,
  Cell,
  Tooltip,
  Legend,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
} from "recharts";

interface ExpenseChartsProps {
  pieData: { name: string; value: number }[];
  barData: { name: string; total: number }[];
}

const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042", "#AA00FF"];

export const Summary: React.FC<ExpenseChartsProps> = ({
  pieData,
  barData,
}) => {
  return (
    <div style={styles.container}>
      <div style={styles.chartBox}>
        <h3 style={styles.chartTitle}>Expenses by Category</h3>
        <PieChart width={300} height={300}>
          <Pie
            data={pieData}
            dataKey="value"
            nameKey="name"
            cx="50%"
            cy="50%"
            outerRadius={100}
            label
          >
            {pieData.map((_entry, index) => (
              <Cell
                key={`cell-${index}`}
                fill={COLORS[index % COLORS.length]}
              />
            ))}
          </Pie>
          <Tooltip />
          <Legend verticalAlign="bottom" height={36} />
        </PieChart>
      </div>

      <div style={styles.chartBox}>
        <h3 style={styles.chartTitle}>Expenses Over Time</h3>
        <BarChart width={500} height={300} data={barData}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="name" tick={{ fill: "#666" }} />
          <YAxis tick={{ fill: "#666" }} />
          <Tooltip />
          <Legend verticalAlign="bottom" height={36} />
          <Bar dataKey="total" fill="#8884d8" radius={[4, 4, 0, 0]} />
        </BarChart>
      </div>
    </div>
  );
};

const styles: { [key: string]: React.CSSProperties } = {
  container: {
    display: "flex",
    justifyContent: "center",
    gap: "3rem",
    padding: "2rem",
    background: "#fff",
    borderRadius: "12px",
    boxShadow: "0 4px 15px rgba(0,0,0,0.1)",
    flexWrap: "wrap",
  },
  chartBox: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    background: "#f0f4f8",
    padding: "1.5rem",
    borderRadius: "12px",
    boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
    minWidth: "450px",
  },
  chartTitle: {
    marginBottom: "1rem",
    fontFamily: "'Segoe UI', Tahoma, Geneva, Verdana, sans-serif",
    color: "#333",
  },
};

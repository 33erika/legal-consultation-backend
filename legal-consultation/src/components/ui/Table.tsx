import { type ReactNode } from 'react';
import { cn } from '@/lib/utils';

export interface TableProps {
  children: ReactNode;
  className?: string;
}

export function Table({ children, className }: TableProps) {
  return (
    <div className={cn('overflow-x-auto', className)}>
      <table className="min-w-full divide-y divide-gray-200">{children}</table>
    </div>
  );
}

export function TableHeader({ children, className }: { children: ReactNode; className?: string }) {
  return <thead className={cn('bg-gray-50', className)}>{children}</thead>;
}

export function TableBody({ children, className }: { children: ReactNode; className?: string }) {
  return <tbody className={cn('bg-white divide-y divide-gray-200', className)}>{children}</tbody>;
}

export interface TableRowProps {
  children: ReactNode;
  className?: string;
  onClick?: () => void;
}

export function TableRow({ children, className, onClick }: TableRowProps) {
  return (
    <tr
      className={cn(
        onClick && 'cursor-pointer hover:bg-gray-50',
        className
      )}
      onClick={onClick}
    >
      {children}
    </tr>
  );
}

export interface TableCellProps {
  children: ReactNode;
  className?: string;
}

export function TableCell({ children, className }: TableCellProps) {
  return <td className={cn('px-6 py-4 whitespace-nowrap text-sm text-gray-900', className)}>{children}</td>;
}

export function TableHead({ children, className }: { children: ReactNode; className?: string }) {
  return (
    <th
      scope="col"
      className={cn('px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider', className)}
    >
      {children}
    </th>
  );
}

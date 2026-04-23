import { type ReactNode } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import {
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
  ChevronRight as RowChevron,
  Inbox,
} from "lucide-react";

export interface Column<T> {
  header: string;
  accessor: keyof T | ((row: T) => ReactNode);
  className?: string;
}

interface DataTableProps<T> {
  columns: Column<T>[];
  data: T[];
  total?: number;
  page?: number;
  pageSize?: number;
  onPageChange?: (page: number) => void;
  onRowClick?: (row: T) => void;
  keyExtractor: (row: T) => string;
  emptyMessage?: string;
  emptyDescription?: string;
  isLoading?: boolean;
}

function Skeleton({ className }: { className?: string }) {
  return (
    <div
      className={`animate-pulse rounded bg-muted/70 ${className ?? "h-4 w-full"}`}
    />
  );
}

export function DataTable<T>({
  columns,
  data,
  total = 0,
  page = 1,
  pageSize = 20,
  onPageChange,
  onRowClick,
  keyExtractor,
  emptyMessage = "No results found",
  emptyDescription,
  isLoading = false,
}: DataTableProps<T>) {
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const showPagination = onPageChange && total > pageSize;
  const rangeStart = (page - 1) * pageSize + 1;
  const rangeEnd = Math.min(page * pageSize, total);

  // Generate visible page numbers (show up to 5 pages centered around current)
  const getPageNumbers = (): (number | "ellipsis")[] => {
    if (totalPages <= 5) {
      return Array.from({ length: totalPages }, (_, i) => i + 1);
    }
    const pages: (number | "ellipsis")[] = [];
    if (page <= 3) {
      pages.push(1, 2, 3, 4, "ellipsis", totalPages);
    } else if (page >= totalPages - 2) {
      pages.push(
        1,
        "ellipsis",
        totalPages - 3,
        totalPages - 2,
        totalPages - 1,
        totalPages
      );
    } else {
      pages.push(1, "ellipsis", page - 1, page, page + 1, "ellipsis", totalPages);
    }
    return pages;
  };

  return (
    <div className="space-y-3">
      <div className="rounded-lg border bg-card shadow-sm overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow className="hover:bg-transparent">
              {columns.map((col) => (
                <TableHead key={col.header} className={col.className}>
                  {col.header}
                </TableHead>
              ))}
              {onRowClick && <TableHead className="w-8" />}
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              Array.from({ length: Math.min(pageSize, 5) }).map((_, i) => (
                <TableRow key={`skeleton-${i}`} className="hover:bg-transparent">
                  {columns.map((col) => (
                    <TableCell key={col.header} className={col.className}>
                      <Skeleton
                        className={`h-4 ${
                          col.className?.includes("text-right")
                            ? "ml-auto w-16"
                            : "w-[60%]"
                        }`}
                      />
                    </TableCell>
                  ))}
                  {onRowClick && (
                    <TableCell>
                      <Skeleton className="size-4" />
                    </TableCell>
                  )}
                </TableRow>
              ))
            ) : data.length === 0 ? (
              <TableRow className="hover:bg-transparent">
                <TableCell
                  colSpan={columns.length + (onRowClick ? 1 : 0)}
                  className="h-40"
                >
                  <div className="flex flex-col items-center justify-center text-center py-4">
                    <div className="flex size-11 items-center justify-center rounded-xl bg-muted mb-2.5">
                      <Inbox className="size-5 text-muted-foreground/50" />
                    </div>
                    <p className="text-sm font-medium text-muted-foreground">
                      {emptyMessage}
                    </p>
                    {emptyDescription && (
                      <p className="text-xs text-muted-foreground/60 mt-0.5">
                        {emptyDescription}
                      </p>
                    )}
                  </div>
                </TableCell>
              </TableRow>
            ) : (
              data.map((row) => (
                <TableRow
                  key={keyExtractor(row)}
                  className={
                    onRowClick
                      ? "cursor-pointer group/row"
                      : undefined
                  }
                  onClick={() => onRowClick?.(row)}
                >
                  {columns.map((col) => (
                    <TableCell key={col.header} className={col.className}>
                      {typeof col.accessor === "function"
                        ? col.accessor(row)
                        : (row[col.accessor] as ReactNode)}
                    </TableCell>
                  ))}
                  {onRowClick && (
                    <TableCell className="w-8 pr-3">
                      <RowChevron className="size-4 text-muted-foreground/30 transition-all group-hover/row:text-muted-foreground/70 group-hover/row:translate-x-0.5" />
                    </TableCell>
                  )}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {showPagination && (
        <div className="flex items-center justify-between">
          <p className="text-xs text-muted-foreground tabular-nums">
            <span className="font-medium text-foreground/80">
              {rangeStart}&ndash;{rangeEnd}
            </span>{" "}
            of {total.toLocaleString()} results
          </p>
          <div className="flex items-center gap-1">
            <Button
              variant="ghost"
              size="icon"
              className="size-8"
              onClick={() => onPageChange(1)}
              disabled={page <= 1}
            >
              <ChevronsLeft className="size-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="size-8"
              onClick={() => onPageChange(page - 1)}
              disabled={page <= 1}
            >
              <ChevronLeft className="size-4" />
            </Button>

            <div className="flex items-center gap-0.5 mx-1">
              {getPageNumbers().map((p, i) =>
                p === "ellipsis" ? (
                  <span
                    key={`ellipsis-${i}`}
                    className="w-8 text-center text-xs text-muted-foreground/50 select-none"
                  >
                    ...
                  </span>
                ) : (
                  <Button
                    key={p}
                    variant={p === page ? "default" : "ghost"}
                    size="icon"
                    className={`size-8 text-xs font-medium tabular-nums ${
                      p === page ? "" : "text-muted-foreground"
                    }`}
                    onClick={() => onPageChange(p)}
                  >
                    {p}
                  </Button>
                )
              )}
            </div>

            <Button
              variant="ghost"
              size="icon"
              className="size-8"
              onClick={() => onPageChange(page + 1)}
              disabled={page >= totalPages}
            >
              <ChevronRight className="size-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="size-8"
              onClick={() => onPageChange(totalPages)}
              disabled={page >= totalPages}
            >
              <ChevronsRight className="size-4" />
            </Button>
          </div>
        </div>
      )}

      {!showPagination && data.length > 0 && (
        <p className="text-xs text-muted-foreground">
          {total.toLocaleString()} {total === 1 ? "result" : "results"}
        </p>
      )}
    </div>
  );
}

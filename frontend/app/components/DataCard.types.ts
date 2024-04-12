export interface DataCardProps {
    data: {
      table_name?: string;
      description?: string;
      database_name?: string;
      column_name?: string;
      health?: string;
      data_classification?: string;
    };
  }
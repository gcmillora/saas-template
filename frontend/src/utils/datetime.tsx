import { DateTime } from "luxon";
export const formatDateFromString = (
  dateStr: string | undefined | null,
  strIfEmpty?: string,
) => {
  if (!dateStr) return strIfEmpty ?? "";

  const date = DateTime.fromISO(dateStr);

  if (!date.isValid) {
    return dateStr;
  }

  return getFormatString(date);
};

const getFormatString = (date: DateTime): string => {
  // const timezone = DateTime.local().zoneName;
  const format = "yyyy/MM/dd HH:mm";

  const result = date.toFormat(format);

  return result;
};

import "./App.css";
import {
  ChangeEvent,
  FormEvent,
  MutableRefObject,
  ReactNode,
  useEffect,
  useRef,
  useState,
} from "react";
import { main } from "../wailsjs/go/models";
import * as runtime from "../wailsjs/runtime/runtime";
import * as app from "../wailsjs/go/main/App";

export default function App() {
  const matrix = useRef<number[][]>([]);
  const [sudoku, setSudoku] = useState<main.Sudoku>({
    size: 0,
    values: [],
    locked: [],
  });
  const [startFromOne, setStartFromOne] = useState(false);

  function updateSudoku(s: main.Sudoku) {
    setSudoku(s);
    matrix.current = s.values;
    const p = new Array(s.size);

    for (let i = 0; i < s.size; i++) {
      p[i] = new Array(s.size);
      for (let j = 0; j < s.size; j++) {
        p[i][j] = [];
      }
    }
  }

  useEffect(() => {
    runtime.EventsEmit("request_possibles");
  }, []);

  app.GetSudoku().then((s) => {
    if (!s) return;

    updateSudoku(s);
  });

  async function createSudoku(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = new FormData(e.currentTarget);

    const size = parseInt((form.get("size") as string) ?? "0");
    await app.NewSudoku(size);
    const s = await app.GetSudoku();
    updateSudoku(s);
  }

  async function lockSudoku() {
    const locks = await app.LockCells();
    setSudoku({ size: sudoku.size, values: sudoku.values, locked: locks });
  }

  async function unlockSudoku() {
    const locks = await app.UnlockCells();
    setSudoku({ size: sudoku.size, values: sudoku.values, locked: locks });
  }

  return (
    <div className="flex flex-col items-center justify-center h-full">
      <div className="flex items-end">
        <form className="flex flex-col w-fit p-4 gap-2" onSubmit={createSudoku}>
          Size:
          <input
            className="border bg-neutral-100 px-4 block w-20"
            type="number"
            name="size"
            placeholder={"9"}
          />
          <button className="bg-neutral-800 text-neutral-100 py-1 rounded w-20 hover:bg-neutral-700 active:bg-neutral-950">
            Create
          </button>
        </form>

        <div className="flex flex-col justify-end items-start p-4 gap-2">
          <div className="flex gap-3">
            <label>From 1:</label>
            <input
              type="checkbox"
              onChange={() => setStartFromOne(!startFromOne)}
              className=""
            />
          </div>
          <button
            className="bg-neutral-800 text-neutral-100 py-1 rounded w-20 hover:bg-neutral-700 active:bg-neutral-950 h-fit"
            onClick={lockSudoku}
          >
            Lock
          </button>
          <button
            className="bg-neutral-800 text-neutral-100 py-1 rounded w-20 hover:bg-neutral-700 active:bg-neutral-950 h-fit"
            onClick={unlockSudoku}
          >
            Unlock
          </button>
        </div>
      </div>

      <Sudoku matrix={matrix} sudoku={sudoku} startFromOne={startFromOne} />
    </div>
  );
}

type SudokuProps = {
  matrix: MutableRefObject<number[][]>;
  sudoku: main.Sudoku;
  startFromOne: boolean;
};

function Sudoku({ sudoku, matrix, startFromOne }: SudokuProps) {
  const fields: ReactNode[] = new Array(sudoku.size * sudoku.size);

  for (let i = 0; i < sudoku.size; i++) {
    for (let j = 0; j < sudoku.size; j++) {
      fields[i * sudoku.size + j] = (
        <SudokuField
          key={`${sudoku.size}-${i * sudoku.size + j}`}
          i={i}
          j={j}
          value={matrix.current[i][j]}
          sudoku={sudoku}
          startFromOne={startFromOne}
        />
      );
    }
  }

  return (
    <div
      className="grid w-fit gap-[1px] border bg-neutral-200"
      style={{
        gridTemplateColumns: `repeat(${sudoku.size}, 1fr)`,
      }}
    >
      {fields}
    </div>
  );
}

type SudokuField = {
  i: number;
  j: number;
  value: number;
  sudoku: main.Sudoku;
  startFromOne: boolean;
};

function SudokuField({ i, j, value: v, sudoku, startFromOne }: SudokuField) {
  const [value, setValue] = useState(v);
  const [valid, setValid] = useState(true);
  const [possibles, setPossibles] = useState<number[]>([]);

  function invalidFieldListener() {
    setValid(false);
  }

  useEffect(() => {
    runtime.EventsOn(
      `invalid_field:${sudoku.size}-${i}:${j}`,
      invalidFieldListener,
    );

    return () => runtime.EventsOff(`invalid_field:${sudoku.size}-${i}:${j}`);
  }, []);

  function onPossibilityUpdate(possible: number[]) {
    setPossibles(possible ? possible : []);
  }

  useEffect(() => {
    runtime.EventsOn(`possibility_update:${i}-${j}`, onPossibilityUpdate);
    return () => {
      runtime.EventsOff(`possibility_update:${i}-${j}`);
    };
  }, []);

  async function onChange(e: ChangeEvent<HTMLInputElement>) {
    const s = e.currentTarget.value;
    let v: number;

    if (s.length === 0) {
      v = -1;
    } else {
      v = parseInt(e.currentTarget.value, 16);
      if (startFromOne) {
        v -= 1;
      }
    }

    if (v >= sudoku.size || v < -1 || Number.isNaN(v)) return;

    setValid(true);
    try {
      await app.InitCell(i, j, v);
    } catch (e) {
      console.error(e);
      return;
    }
    setValue(v);
  }

  return (
    <div
      className={`relative bg-white overflow-hidden w-10 h-10 ${border(i, j, sudoku.size)}`}
    >
      <input
        value={displayValue(value, startFromOne)}
        className={`block w-full h-full text-center disabled:bg-neutral-200 ${value != -1 ? "bg-blue-100" : ""} ${!valid ? "bg-red-500" : ""}`}
        type="text"
        onChange={onChange}
        disabled={sudoku.locked[i][j]}
      />
      <p className="text-[9px] absolute right-[1px] bottom-0">
        {possibles.map((x) =>
          (startFromOne ? x + 1 : x).toString(16).toUpperCase(),
        )}
      </p>
    </div>
  );
}

function displayValue(value: number, fromOne: boolean) {
  if (value === -1 || value === undefined) return "";
  return (fromOne ? value + 1 : value).toString(16).toUpperCase();
}

function border(i: number, j: number, size: number) {
  const block = Math.floor(Math.sqrt(size));

  let style = ["border-black"];
  if ((i % block) + 1 == block && i < size - 1) {
    style.push("border-b");
  }

  if ((j % block) + 1 == block && j < size - 1) {
    style.push("border-r");
  }

  return style.join(" ");
}

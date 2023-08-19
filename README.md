# PavoSQL
A simple SQL database written in pure Go.

---

## B-Tree Page Format

### B-Tree Page Header

| Offset| Size              | Description
|-------|-------------------|------------
| 0     | 1                 | Type: first byte defines the type of the page
| 1     | 2                 | nCells: uint16 number defining the number of cells currently stored on the page
| 3     | 2                 | cellSize: uint16 number defining the size of a single cell that is allowed to be stored on the page (only cells of the same type and size are stored on a single page)
| 5     | nCells*cellSize   | Cells: a list of nCells byte chunks the size of cellSize representing the actual cells stored on the page

### Cell Format

| Offset    | Size  | Description
|-----------|-------|------------
| 0         | 1     | Type: first byte defines the type of the cell
| 1         | 2     | kSize: uint16 number defining the size of the key stored in the cell
| 3         | 2     | vSize: uint16 number defining the size of the value stored in the cell
| 5         | kSize | Key: bytes that represent the key stored in the cell
| 5+kSize   | vSize | Val: bytes that represent the value stored in the cell

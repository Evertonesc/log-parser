# Quake Log Parser

A log parser written in Go that collects and reports kill information for every match present in a log file of Quake III Arena.

This parser was built with concurrency, gathering the relevant log lines with the desired information and passing them
to a CoR (Chain of Responsibility) to process the log lines in the handlers. The implementation of the Chain of Responsibility Design Pattern 
makes the main functionality of the parser easy to understand and flexible enough to increase its features.

---

## Execution

**Run with Go**

``go run main.go``

**Run with Makefile**

``make exec``




---

## Report Outputs

### Match Kill Data

```json
{
   "game_21":{
      "total_kills":131,
      "players":[
         "Isgalamido",
         "Oootsimo",
         "Dono da Bola",
         "Assasinu Credi",
         "Zeh",
         "Mal"
      ],
      "kills":{
         "Assasinu Credi":16,
         "Dono da Bola":12,
         "Isgalamido":17,
         "Mal":6,
         "Oootsimo":21,
         "Zeh":19
      }
   }
}
```

### Deaths grouped by Death Cause

```json
    {
        "game_21": {
            "kills_by_means": {
                "MOD_FALLING": 3,
                "MOD_MACHINEGUN": 4,
                "MOD_RAILGUN": 9,
                "MOD_ROCKET": 37,
                "MOD_ROCKET_SPLASH": 60,
                "MOD_SHOTGUN": 4,
                "MOD_TRIGGER_HURT": 14
            }
        }
    }
```

---


### Complete Report Sample

_Matches Report - 16/07/2024 16:45_
```json
[
   {
      "game_21":{
         "total_kills":131,
         "players":[
            "Isgalamido",
            "Oootsimo",
            "Dono da Bola",
            "Assasinu Credi",
            "Zeh",
            "Mal"
         ],
         "kills":{
            "Assasinu Credi":16,
            "Dono da Bola":12,
            "Isgalamido":17,
            "Mal":6,
            "Oootsimo":21,
            "Zeh":19
         }
      }
   }
]
```
_Deaths by Death cause - 16/07/2024 16:45_

```json
[
   {
      "game_21":{
         "kills_by_means":{
            "MOD_FALLING":3,
            "MOD_MACHINEGUN":4,
            "MOD_RAILGUN":9,
            "MOD_ROCKET":37,
            "MOD_ROCKET_SPLASH":60,
            "MOD_SHOTGUN":4,
            "MOD_TRIGGER_HURT":14
         }
      }
   }
]
```

reports generated in 11 ms
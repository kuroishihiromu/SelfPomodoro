import torch
import torch.nn as nn
from torch.nn import GRU, Linear, Sequential, ReLU, Dropout

class FocusScoreModel(nn.Module):
    """集中度スコアを予測するモデル"""

    def __init__(self, input_size: int, hidden_size: int):
        super(FocusScoreModel, self).__init__()

        self.gru1 = nn.GRU(
            input_size=input_size,
            hidden_size=hidden_size,
            batch_first=True,
        )

        self.gru2 = nn.GRU(
            input_size=hidden_size,
            hidden_size=hidden_size,
            batch_first=False,
        )

        self.fc1 = nn.Linear(hidden_size, 25)

        self.fc2 = nn.Linear(25, 1)
    

    def forward(self, x):
        gru1_out = self.gru1(x)

        gru2_out, _ = self.gru2(gru1_out)

        last_step_out = gru2_out[:, -1, :]

        x = torch.relu(self.fc1(last_step_out))

        x = self.fc2(x)


        return x
